package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dburton90/wrt/internal/gitutil"
	"github.com/dburton90/wrt/internal/registry"
	"github.com/dburton90/wrt/internal/task"
	"github.com/spf13/cobra"
)

var closeCmd = &cobra.Command{
	Use:               "close <task-name>",
	Short:             "Close a task: save patches and remove worktrees",
	Args:              cobra.ExactArgs(1),
	RunE:              runClose,
	ValidArgsFunction: completeOpenTaskNames,
}

func init() {
	rootCmd.AddCommand(closeCmd)
}

func runClose(_ *cobra.Command, args []string) error {
	_, taskRoot := mustTaskRoot()
	name := args[0]

	taskDir, open, err := task.Find(taskRoot, name)
	if err != nil {
		return err
	}
	if !open {
		return fmt.Errorf("task %q is already closed", name)
	}

	t, err := task.Load(taskDir)
	if err != nil {
		return err
	}

	// Save patches and remove worktrees for each repo
	for repoName, rs := range t.Repositories {
		r, err := registry.Load(taskRoot, repoName)
		if err != nil {
			return fmt.Errorf("cannot close task %q: repo %q listed in task.json is not in the registry. Restore the registry entry or remove %q from task.json before closing", name, repoName, repoName)
		}

		repoDir := filepath.Join(taskDir, "repositories", repoName)
		codeWorktree := filepath.Join(repoDir, "code")

		// Format patch for main worktree
		if _, err := os.Stat(codeWorktree); err == nil {
			patchFile := filepath.Join(repoDir, repoName+".patch")
			fmt.Printf("Saving patch for %s...\n", repoName)
			if err := gitutil.FormatPatch(codeWorktree, rs.BaseBranch, patchFile); err != nil {
				return fmt.Errorf("saving patch for %s: %w", repoName, err)
			}
			fmt.Printf("Removing worktree for %s...\n", repoName)
			if err := gitutil.WorktreeRemove(r.Path, codeWorktree); err != nil {
				return fmt.Errorf("removing worktree for %s: %w", repoName, err)
			}
		}

		// Format patch and remove worktrees for backports
		for _, bp := range rs.BackportBranches {
			bpWorktree := filepath.Join(repoDir, "backports", bp.Version)
			if _, err := os.Stat(bpWorktree); err == nil {
				patchFile := filepath.Join(repoDir, repoName+"-backport-"+bp.Version+".patch")
				fmt.Printf("Saving patch for %s backport %s...\n", repoName, bp.Version)
				if err := gitutil.FormatPatch(bpWorktree, bp.BaseBranch, patchFile); err != nil {
					return fmt.Errorf("saving backport patch for %s/%s: %w", repoName, bp.Version, err)
				}
				fmt.Printf("Removing backport worktree %s/%s...\n", repoName, bp.Version)
				if err := gitutil.WorktreeRemove(r.Path, bpWorktree); err != nil {
					return fmt.Errorf("removing backport worktree for %s/%s: %w", repoName, bp.Version, err)
				}
			}
		}
	}

	// Move task to closed
	closedDir := task.ClosedDir(taskRoot, name)
	if err := os.MkdirAll(filepath.Dir(closedDir), 0o755); err != nil {
		return fmt.Errorf("creating closed dir: %w", err)
	}
	if err := os.Rename(taskDir, closedDir); err != nil {
		return fmt.Errorf("moving task to closed: %w", err)
	}

	fmt.Printf("Task %q closed. Patches saved to %s\n", name, filepath.Join(closedDir, "repositories"))
	return nil
}
