package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/dburton90/wrt/internal/gitutil"
	"github.com/dburton90/wrt/internal/registry"
	"github.com/dburton90/wrt/internal/task"
	"github.com/dburton90/wrt/internal/taskctx"
	"github.com/spf13/cobra"
)

var backportCmd = &cobra.Command{
	Use:   "backport",
	Short: "Manage backport worktrees",
}

var backportAddCmd = &cobra.Command{
	Use:   "add <repo-name> <version>",
	Short: "Add a backport worktree to the current task",
	Args:  cobra.ExactArgs(2),
	RunE:  runBackportAdd,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completeBackportArgs(args, toComplete)
	},
}

func init() {
	rootCmd.AddCommand(backportCmd)
	backportCmd.AddCommand(backportAddCmd)
}

func runBackportAdd(_ *cobra.Command, args []string) error {
	cfg, taskRoot := mustTaskRoot()
	repoName := args[0]
	version := args[1]

	taskDir, err := taskctx.Find()
	if err != nil {
		return err
	}

	t, err := task.Load(taskDir)
	if err != nil {
		return err
	}

	rs, exists := t.Repositories[repoName]
	if !exists {
		return fmt.Errorf("repo %q is not part of this task\n\nRun `wrt repo add %s` first", repoName, repoName)
	}

	// Check not already added
	for _, bp := range rs.BackportBranches {
		if bp.Version == version {
			return fmt.Errorf("backport %q for repo %q already exists in this task", version, repoName)
		}
	}

	r, err := registry.Load(taskRoot, repoName)
	if err != nil {
		return err
	}

	baseBranch, ok := r.BackportBranches[version]
	if !ok {
		return fmt.Errorf(
			"version %q is not configured for repo %q\n\nRun `wrt repo list` to see available backport versions",
			version, repoName,
		)
	}

	branch := task.ResolveBranch(r.BackportBranchTemplate, cfg.Username(), t.Name, version)
	worktreePath := filepath.Join(taskDir, "repositories", repoName, "backports", version)

	fmt.Printf("Creating backport branch %s from %s...\n", branch, baseBranch)
	if err := gitutil.WorktreeAdd(r.Path, worktreePath, branch, baseBranch); err != nil {
		return fmt.Errorf("creating backport worktree: %w", err)
	}

	rs.BackportBranches = append(rs.BackportBranches, task.BackportState{
		Version:    version,
		BaseBranch: baseBranch,
		Branch:     branch,
	})
	t.Repositories[repoName] = rs
	if err := task.Save(taskDir, t); err != nil {
		return fmt.Errorf("updating task.json: %w", err)
	}

	fmt.Printf("Added backport %q for repo %q\n  worktree: %s\n  branch:   %s\n", version, repoName, worktreePath, branch)
	return nil
}

func completeBackportArgs(args []string, _ string) ([]string, cobra.ShellCompDirective) {
	taskRoot, err := loadTaskRootForCompletion()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if len(args) == 0 {
		// Complete repo names present in current task
		taskDir, err := taskctx.Find()
		if err != nil {
			return registry.Names(taskRoot), cobra.ShellCompDirectiveNoFileComp
		}
		t, err := task.Load(taskDir)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		names := make([]string, 0, len(t.Repositories))
		for name := range t.Repositories {
			names = append(names, name)
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	}

	if len(args) == 1 {
		// Complete backport versions for the given repo
		repoName := args[0]
		r, err := registry.Load(taskRoot, repoName)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		versions := make([]string, 0, len(r.BackportBranches))
		for v := range r.BackportBranches {
			versions = append(versions, v)
		}
		return versions, cobra.ShellCompDirectiveNoFileComp
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}
