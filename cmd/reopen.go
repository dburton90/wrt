package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dburton90/wrt/internal/gitutil"
	"github.com/dburton90/wrt/internal/task"
	"github.com/spf13/cobra"
)

var reopenCmd = &cobra.Command{
	Use:               "reopen <task-name>",
	Short:             "Reopen a closed task: recreate worktrees and apply patches",
	Args:              cobra.ExactArgs(1),
	RunE:              runReopen,
	ValidArgsFunction: completeClosedTaskNames,
}

func init() {
	rootCmd.AddCommand(reopenCmd)
}

func runReopen(_ *cobra.Command, args []string) error {
	_, taskRoot := mustTaskRoot()
	name := args[0]

	taskDir, open, err := task.Find(taskRoot, name)
	if err != nil {
		return err
	}
	if open {
		return fmt.Errorf("task %q is already open", name)
	}

	t, err := task.Load(taskDir)
	if err != nil {
		return err
	}

	// Move to open first
	openDir := task.OpenDir(taskRoot, name)
	if err := os.MkdirAll(filepath.Dir(openDir), 0o755); err != nil {
		return fmt.Errorf("creating open dir: %w", err)
	}
	if err := os.Rename(taskDir, openDir); err != nil {
		return fmt.Errorf("moving task to open: %w", err)
	}

	var conflicted bool

	for repoName, rs := range t.Repositories {
		repoDir := filepath.Join(openDir, "repositories", repoName)
		codeWorktree := filepath.Join(repoDir, "code")

		// Recreate main worktree (branch already exists from before close)
		fmt.Printf("Recreating worktree for %s on %s...\n", repoName, rs.TaskBranch)
		if err := gitutil.WorktreeAddExisting(rs.RepoPath, codeWorktree, rs.TaskBranch); err != nil {
			return fmt.Errorf("recreating worktree for %s: %w", repoName, err)
		}

		// Apply patch only if branch has no commits ahead of base (branch was recreated fresh)
		patchFile := filepath.Join(repoDir, repoName+".patch")
		if gitutil.HasCommitsAhead(codeWorktree, rs.BaseBranch) {
			fmt.Printf("  Branch already has commits, skipping patch for %s\n", repoName)
		} else if err := gitutil.ApplyPatch(codeWorktree, patchFile); err != nil {
			var conflictErr *gitutil.PatchConflictError
			if errors.As(err, &conflictErr) {
				conflicted = true
				fmt.Printf("\n⚠ Patch conflict in %s:\n%s\n", repoName, conflictErr.Output)
				fmt.Printf("Resolve conflicts, then run:\n  cd %s && git am --continue\n", codeWorktree)
				fmt.Printf("Or to abandon the patch:\n  cd %s && git am --abort\n\n", codeWorktree)
			} else {
				return fmt.Errorf("applying patch for %s: %w", repoName, err)
			}
		}

		// Recreate backport worktrees (branches already exist)
		for _, bp := range rs.BackportBranches {
			bpWorktree := filepath.Join(repoDir, "backports", bp.Version)
			fmt.Printf("Recreating backport worktree %s/%s...\n", repoName, bp.Version)
			if err := gitutil.WorktreeAddExisting(rs.RepoPath, bpWorktree, bp.Branch); err != nil {
				return fmt.Errorf("recreating backport worktree %s/%s: %w", repoName, bp.Version, err)
			}
			bpPatchFile := filepath.Join(repoDir, repoName+"-backport-"+bp.Version+".patch")
			if gitutil.HasCommitsAhead(bpWorktree, bp.BaseBranch) {
				fmt.Printf("  Branch already has commits, skipping patch for %s/%s\n", repoName, bp.Version)
			} else if err := gitutil.ApplyPatch(bpWorktree, bpPatchFile); err != nil {
				var conflictErr *gitutil.PatchConflictError
				if errors.As(err, &conflictErr) {
					conflicted = true
					fmt.Printf("\n⚠ Patch conflict in %s backport %s:\n%s\n", repoName, bp.Version, conflictErr.Output)
					fmt.Printf("Resolve conflicts, then run:\n  cd %s && git am --continue\n\n", bpWorktree)
				} else {
					return fmt.Errorf("applying backport patch for %s/%s: %w", repoName, bp.Version, err)
				}
			}
		}
	}

	if conflicted {
		fmt.Printf("Task %q reopened with conflicts. Resolve them manually before continuing.\n", name)
		os.Exit(1)
	}

	fmt.Printf("Task %q reopened at %s\n", name, openDir)
	return nil
}
