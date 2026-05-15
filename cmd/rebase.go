package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/dburton90/wrt/internal/gitutil"
	"github.com/dburton90/wrt/internal/registry"
	"github.com/dburton90/wrt/internal/task"
	"github.com/dburton90/wrt/internal/taskctx"
	"github.com/spf13/cobra"
)

var rebaseCmd = &cobra.Command{
	Use:   "rebase",
	Short: "Rebase every worktree in the current task onto its base branch's upstream",
	Args:  cobra.NoArgs,
	RunE:  runRebase,
}

func init() {
	rootCmd.AddCommand(rebaseCmd)
}

type rebaseTarget struct {
	worktreePath string
	baseBranch   string
	repoPath     string
	label        string
}

type targetResult struct {
	label  string
	base   string
	status string
}

const (
	statusOK         = "ok"
	statusConflict   = "CONFLICT"
	statusNoUpstream = "SKIPPED — no upstream"
	statusFetchFail  = "FETCH FAILED"
	statusPending    = "pending"
)

func runRebase(_ *cobra.Command, _ []string) error {
	_, taskRoot := mustTaskRoot()

	taskDir, err := taskctx.Find()
	if err != nil {
		return err
	}
	t, err := task.Load(taskDir)
	if err != nil {
		return err
	}

	targets, err := buildRebaseTargets(taskRoot, taskDir, t)
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		fmt.Println("No worktrees to rebase.")
		return nil
	}

	if err := preflightClean(targets); err != nil {
		return err
	}

	results := processTargets(targets)
	printRebaseSummary(t.Name, results)

	for _, r := range results {
		if r.status == statusConflict || r.status == statusFetchFail {
			os.Exit(1)
		}
	}
	return nil
}

func buildRebaseTargets(taskRoot, taskDir string, t *task.Task) ([]rebaseTarget, error) {
	repoNames := make([]string, 0, len(t.Repositories))
	for name := range t.Repositories {
		repoNames = append(repoNames, name)
	}
	sort.Strings(repoNames)

	var targets []rebaseTarget
	for _, repoName := range repoNames {
		rs := t.Repositories[repoName]
		r, err := registry.Load(taskRoot, repoName)
		if err != nil {
			return nil, fmt.Errorf("cannot rebase task: repo %q listed in task.json is not in the registry. Restore the registry entry or remove %q from task.json", repoName, repoName)
		}
		repoDir := filepath.Join(taskDir, "repositories", repoName)

		targets = append(targets, rebaseTarget{
			worktreePath: filepath.Join(repoDir, "code"),
			baseBranch:   rs.BaseBranch,
			repoPath:     r.Path,
			label:        repoName + "/code",
		})

		bps := append([]task.BackportState{}, rs.BackportBranches...)
		sort.Slice(bps, func(i, j int) bool { return bps[i].Version < bps[j].Version })
		for _, bp := range bps {
			targets = append(targets, rebaseTarget{
				worktreePath: filepath.Join(repoDir, "backports", bp.Version),
				baseBranch:   bp.BaseBranch,
				repoPath:     r.Path,
				label:        repoName + "/backports/" + bp.Version,
			})
		}
	}
	return targets, nil
}

func preflightClean(targets []rebaseTarget) error {
	var dirty []string
	for _, tg := range targets {
		clean, files, err := gitutil.IsClean(tg.worktreePath)
		if err != nil {
			return fmt.Errorf("checking worktree %s: %w", tg.label, err)
		}
		if !clean {
			dirty = append(dirty, fmt.Sprintf("  %s  (%d uncommitted file(s))", tg.label, len(files)))
		}
	}
	if len(dirty) > 0 {
		fmt.Fprintln(os.Stderr, "Cannot rebase: the following worktrees have uncommitted changes:")
		for _, d := range dirty {
			fmt.Fprintln(os.Stderr, d)
		}
		fmt.Fprintln(os.Stderr, "Commit or stash, then re-run `wrt rebase`.")
		os.Exit(1)
	}
	return nil
}

func processTargets(targets []rebaseTarget) []targetResult {
	results := make([]targetResult, len(targets))
	for i, tg := range targets {
		results[i].label = tg.label

		remote, ref, err := gitutil.Upstream(tg.repoPath, tg.baseBranch)
		if err != nil {
			if errors.Is(err, gitutil.ErrNoUpstream) {
				results[i].base = "?"
				results[i].status = statusNoUpstream
				fmt.Printf("  %s (← ?)  SKIPPED — base branch %q has no upstream\n", tg.label, tg.baseBranch)
				continue
			}
			results[i].base = "?"
			results[i].status = statusFetchFail
			fmt.Printf("  %s  upstream lookup failed: %v\n", tg.label, err)
			markPending(results, i+1, targets)
			return results
		}
		base := remote + "/" + ref
		results[i].base = base

		fmt.Printf("  %s (← %s)  fetching %s...\n", tg.label, base, remote)
		if err := gitutil.Fetch(tg.repoPath, remote); err != nil {
			results[i].status = statusFetchFail
			fmt.Printf("  fetch failed: %v\n", err)
			markPending(results, i+1, targets)
			return results
		}

		fmt.Printf("  %s rebasing onto %s...\n", tg.label, base)
		if err := gitutil.Rebase(tg.worktreePath, base); err != nil {
			results[i].status = statusConflict
			files, _ := gitutil.ConflictedFiles(tg.worktreePath)
			fmt.Printf("\n⚠ Conflict in %s\n", tg.label)
			if len(files) > 0 {
				fmt.Println("  files:")
				for _, f := range files {
					fmt.Printf("    %s\n", f)
				}
			}
			fmt.Printf("  Resolve, then run:\n    cd %s\n    git rebase --continue\n  Then re-run `wrt rebase` to continue with remaining worktrees.\n\n", tg.worktreePath)
			markPending(results, i+1, targets)
			return results
		}

		results[i].status = statusOK
		fmt.Printf("  %s  ok\n", tg.label)
	}
	return results
}

func markPending(results []targetResult, from int, targets []rebaseTarget) {
	for j := from; j < len(targets); j++ {
		results[j].label = targets[j].label
		results[j].status = statusPending
	}
}

func printRebaseSummary(taskName string, results []targetResult) {
	fmt.Printf("\nRebase summary for task %s:\n", taskName)
	for _, r := range results {
		base := r.base
		if base == "" {
			base = "?"
		}
		fmt.Printf("  %-50s  %s\n", r.label+" (← "+base+")", r.status)
	}
}
