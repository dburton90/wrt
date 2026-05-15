package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dburton90/wrt/internal/gitutil"
	"github.com/dburton90/wrt/internal/registry"
	"github.com/dburton90/wrt/internal/task"
	"github.com/dburton90/wrt/internal/taskctx"
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage registered repositories",
}

// repo create
var repoCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Register a new repository",
	Args:  cobra.ExactArgs(1),
	RunE:  runRepoCreate,
}

var repoCreateFlags struct {
	path                   string
	baseBranch             string
	link                   string
	taskBranchTemplate     string
	backportBranchTemplate string
}

// repo list
var repoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered repositories",
	Args:  cobra.NoArgs,
	RunE:  runRepoList,
}

// repo add
var repoAddCmd = &cobra.Command{
	Use:               "add <repo-name>",
	Short:             "Add a registered repository to the current task",
	Args:              cobra.ExactArgs(1),
	RunE:              runRepoAdd,
	ValidArgsFunction: completeRegisteredRepoNames,
}

func init() {
	rootCmd.AddCommand(repoCmd)
	repoCmd.AddCommand(repoCreateCmd, repoListCmd, repoAddCmd)

	repoCreateCmd.Flags().StringVar(&repoCreateFlags.path, "path", "", "Absolute path to the git repository")
	repoCreateCmd.Flags().StringVar(&repoCreateFlags.baseBranch, "base-branch", "", "Default base branch")
	repoCreateCmd.Flags().StringVar(&repoCreateFlags.link, "link", "", "URL to remote (optional)")
	repoCreateCmd.Flags().StringVar(&repoCreateFlags.taskBranchTemplate, "task-template", "", "Task branch template (default: users/{user}/{task-id})")
	repoCreateCmd.Flags().StringVar(&repoCreateFlags.backportBranchTemplate, "backport-template", "", "Backport branch template")
}

func runRepoCreate(cmd *cobra.Command, args []string) error {
	_, taskRoot := mustTaskRoot()
	name := args[0]

	// Error if already exists
	if _, err := os.Stat(registry.RepoPath(taskRoot, name)); err == nil {
		return fmt.Errorf(
			"repo %q already exists\n\nRun `wrt repo list` to see it or edit %s directly",
			name, registry.RepoPath(taskRoot, name),
		)
	}

	// Flag mode: --path was provided, skip all interactive prompts.
	flagMode := cmd.Flags().Changed("path")

	path := repoCreateFlags.path
	if path == "" {
		var err error
		path, err = promptString("Repository path", "", true)
		if err != nil {
			return err
		}
	}

	baseBranch := repoCreateFlags.baseBranch
	if baseBranch == "" {
		if flagMode {
			baseBranch = "main"
		} else {
			var err error
			baseBranch, err = promptString("Default base branch", "main", true)
			if err != nil {
				return err
			}
		}
	}

	link := repoCreateFlags.link
	if link == "" && !flagMode {
		var err error
		link, err = promptOptional("Remote URL")
		if err != nil {
			return err
		}
	}

	taskTemplate := repoCreateFlags.taskBranchTemplate
	if taskTemplate == "" {
		if !flagMode {
			var err error
			taskTemplate, err = promptString("Task branch template", registry.DefaultTaskBranchTemplate, true)
			if err != nil {
				return err
			}
		}
		if taskTemplate == "" {
			taskTemplate = registry.DefaultTaskBranchTemplate
		}
	}

	backportTemplate := repoCreateFlags.backportBranchTemplate
	if backportTemplate == "" {
		if !flagMode {
			var err error
			backportTemplate, err = promptString("Backport branch template", registry.DefaultBackportBranchTemplate, true)
			if err != nil {
				return err
			}
		}
		if backportTemplate == "" {
			backportTemplate = registry.DefaultBackportBranchTemplate
		}
	}

	var backportBranches map[string]string
	if !flagMode {
		var err error
		backportBranches, err = promptBackportBranches()
		if err != nil {
			return err
		}
	}

	r := &registry.Repo{
		Name:                   name,
		Path:                   path,
		Link:                   link,
		DefaultBaseBranch:      baseBranch,
		TaskBranchTemplate:     taskTemplate,
		BackportBranchTemplate: backportTemplate,
		BackportBranches:       backportBranches,
	}

	if err := registry.Save(taskRoot, r); err != nil {
		return fmt.Errorf("saving repo: %w", err)
	}

	fmt.Printf("Repo %q registered at %s\n", name, registry.RepoDir(taskRoot, name))
	return nil
}

func runRepoList(_ *cobra.Command, _ []string) error {
	_, taskRoot := mustTaskRoot()

	repos, err := registry.List(taskRoot)
	if err != nil {
		return err
	}
	if len(repos) == 0 {
		fmt.Println("No repos registered. Run `wrt repo create <name>` to add one.")
		return nil
	}

	fmt.Printf("%-20s  %-12s  %-40s  %s\n", "NAME", "BASE BRANCH", "PATH", "BACKPORTS")
	fmt.Printf("%-20s  %-12s  %-40s  %s\n", "----", "-----------", "----", "---------")
	for _, r := range repos {
		backports := make([]string, 0, len(r.BackportBranches))
		for v := range r.BackportBranches {
			backports = append(backports, v)
		}
		bp := strings.Join(backports, ", ")
		if bp == "" {
			bp = "-"
		}
		path := r.Path
		if len(path) > 38 {
			path = "..." + path[len(path)-35:]
		}
		fmt.Printf("%-20s  %-12s  %-40s  %s\n", r.Name, r.DefaultBaseBranch, path, bp)
	}
	return nil
}

func runRepoAdd(_ *cobra.Command, args []string) error {
	cfg, taskRoot := mustTaskRoot()
	repoName := args[0]

	taskDir, err := taskctx.Find()
	if err != nil {
		return err
	}

	t, err := task.Load(taskDir)
	if err != nil {
		return err
	}

	if _, exists := t.Repositories[repoName]; exists {
		return fmt.Errorf("repo %q is already added to this task", repoName)
	}

	r, err := registry.Load(taskRoot, repoName)
	if err != nil {
		return err
	}

	branch := task.ResolveBranch(r.TaskBranchTemplate, cfg.Username(), t.Name, "")
	worktreePath := filepath.Join(taskDir, "repositories", repoName, "code")

	fmt.Printf("Creating branch %s from %s...\n", branch, r.DefaultBaseBranch)
	if err := gitutil.WorktreeAdd(r.Path, worktreePath, branch, r.DefaultBaseBranch); err != nil {
		return fmt.Errorf("creating worktree: %w", err)
	}

	if t.Repositories == nil {
		t.Repositories = map[string]task.RepoState{}
	}
	t.Repositories[repoName] = task.RepoState{
		BaseBranch:       r.DefaultBaseBranch,
		TaskBranch:       branch,
		BackportBranches: []task.BackportState{},
	}
	if err := task.Save(taskDir, t); err != nil {
		return fmt.Errorf("updating task.json: %w", err)
	}

	fmt.Printf("Added repo %q to task %q\n  worktree: %s\n  branch:   %s\n", repoName, t.Name, worktreePath, branch)
	return nil
}

func completeRegisteredRepoNames(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	taskRoot, err := loadTaskRootForCompletion()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return registry.Names(taskRoot), cobra.ShellCompDirectiveNoFileComp
}
