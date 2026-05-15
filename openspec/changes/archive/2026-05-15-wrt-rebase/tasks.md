## 1. Git helpers

- [x] 1.1 Add `IsClean(worktreePath string) (bool, []string, error)` to `internal/gitutil` — returns whether the worktree is clean and a list of dirty files (parsed from `git status --porcelain`)
- [x] 1.2 Add `Upstream(repoPath, branch string) (remote, ref string, err error)` — runs `git rev-parse --abbrev-ref <branch>@{upstream}`, splits on the first `/`
- [x] 1.3 Add `Fetch(repoPath, remote string) error` — runs `git fetch <remote>` in the base repo
- [x] 1.4 Add `Rebase(worktreePath, ontoRef string) error` — runs `git rebase <onto>` in the worktree
- [x] 1.5 Add helper to enumerate conflicted files post-failure: `ConflictedFiles(worktreePath string) ([]string, error)` via `git diff --name-only --diff-filter=U`

## 2. Command

- [x] 2.1 Create `cmd/rebase.go` with a Cobra command `rebase` (no args)
- [x] 2.2 Detect current task via `taskctx.Find`; load `task.json`
- [x] 2.3 Build the list of worktree targets: for each repo, one entry for `code/` and one per backport
- [x] 2.4 Pre-flight: call `IsClean` for every target; if any are dirty, print the dirty list and exit non-zero — no other side effects
- [x] 2.5 Iterate targets in deterministic order (sorted by repo name; code before backports; backports sorted by version):
    - [x] 2.5.1 Resolve the base branch (from `RepoState.BaseBranch` or `BackportState.BaseBranch`)
    - [x] 2.5.2 Resolve upstream via `Upstream`; on error, record as "no upstream", continue
    - [x] 2.5.3 `Fetch` the resolved remote in the base repo (registry lookup for path)
    - [x] 2.5.4 `Rebase` the worktree onto `<remote>/<ref>`
    - [x] 2.5.5 On non-zero exit: collect conflicted files, print resolution instructions, mark remaining targets as pending, break the loop
    - [x] 2.5.6 On success: print a short line for that target

## 3. Output

- [x] 3.1 Per-target progress line as the loop runs (target name, base, status)
- [x] 3.2 End-of-run summary block listing each target with ok / conflict / no-upstream / pending status
- [x] 3.3 Exit code: 0 if every target ok or no-upstream-skipped; non-zero on conflict or fetch failure

## 4. Wire up

- [x] 4.1 Register `rebaseCmd` under `rootCmd` in `cmd/rebase.go`'s `init()`
- [x] 4.2 Confirm the new command shows up in `wrt --help`

## 5. Manual verification

- [x] 5.1 Run `wrt rebase` in a task with a clean code worktree and an advancing `origin/main` — verify rebase succeeds and summary line is ok
- [x] 5.2 Run with a dirty worktree — verify refusal, no rebase attempted, exit non-zero
- [x] 5.3 Run with a task branch whose base has no upstream — verify "no upstream" report and that other worktrees still complete
- [x] 5.4 Run with a remote named other than `origin` — verify that remote is fetched and the rebase target is `<other>/<branch>`
- [x] 5.5 Induce a conflict (e.g. edit the same file upstream and downstream) — verify stop, conflicted files listed, remaining targets reported pending, and `git rebase --continue` instruction printed
- [x] 5.6 After resolving the conflict and running `wrt rebase` again, verify the remaining targets are processed
