## Why

A task lives on a worktree branch forked from some base (`main`, `release/1.1.1`, …). The base advances on the remote while the task is in flight; the local task branch quietly drifts. Today, catching up means `cd`-ing into every worktree (main + every backport) and running `git fetch && git rebase` by hand. For a multi-repo, multi-backport task that's tedious and error-prone.

`wrt rebase` automates the loop with a strong safety stance: refuse to touch dirty trees, fail loudly on conflicts, walk every worktree of the current task.

## What Changes

- New command `wrt rebase`, operates on the current task (detected via the existing `taskctx.Find`).
- Pre-flight: list every worktree (code + each backport for each repo) and verify each is clean. If any are dirty, the command refuses with the dirty list and changes nothing.
- For each worktree, resolve the upstream of its base branch using git's upstream tracking (`<base>@{upstream}`). Do not assume `origin`.
- `git fetch` the upstream's remote in the base repo (resolved via the repo registry).
- `git rebase <remote>/<branch>` inside the worktree.
- Stop-on-conflict: leave the worktree in a paused rebase state, print resolution instructions, exit non-zero. Remaining worktrees are skipped (so the conflict isn't buried).
- Worktrees whose base branch has no upstream are reported and skipped — they don't block the rest, but they show up clearly in the summary.

## Capabilities

### New Capabilities

- `worktree-rebase`: bring all worktrees of a task up to date with their upstream base branches.

### Modified Capabilities

_(none)_

## Impact

- New command file `cmd/rebase.go`.
- Possibly new helpers in `internal/gitutil` for: working-tree clean check, upstream resolution, fetch, rebase.
- No schema change. No new dependency.
- Depends on `taskctx`, `task`, `registry`, `gitutil` packages — all existing.
- Behavior is per-task only; no global `--all` flag in v1.
