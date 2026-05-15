## Context

A task in flight has 1..N repos, each with a `code/` worktree on the task branch and zero..M `backports/<version>/` worktrees on backport branches. Each of those branches was created from some base (`main`, `release/X.Y`, …). As days pass, those bases advance upstream. Catching up is a per-worktree `git fetch && git rebase`, run from inside the worktree. Multiplied by repos and backports, this is the kind of mechanical loop a tool should own.

The trick is doing it safely. Rebase rewrites history; on a dirty worktree it'd interleave uncommitted changes with the replay; on a conflict it leaves things half-applied. Both states deserve explicit handling.

The other subtle bit: this tool can't assume `origin`. Some users push to `upstream`, some have multiple remotes, some configure tracking explicitly. The right primitive is git's own upstream concept (`<branch>@{upstream}`), which encodes "the canonical remote ref for this branch" regardless of the remote's name.

## Goals / Non-Goals

**Goals:**
- `wrt rebase` brings every worktree of the current task up to date with its base branch's upstream.
- Refuse the entire operation if *any* worktree is dirty — atomic decision, no partial progress before that point.
- Conflict in one worktree halts the whole run; user resolves, then re-runs.
- No assumptions about remote names — use git's upstream tracking.

**Non-Goals:**
- Global flags like `--all` (rebase every open task). Save for later.
- Merge mode, fast-forward-only mode, or rebase strategies (`--rebase-merges`, `--keep-empty`, etc.).
- Auto-resolving conflicts.
- Stashing on the user's behalf.
- Refreshing the base branch in the base repo itself (e.g. fast-forwarding the local `main`). We rebase onto the *remote* ref, so the local base branch's state is irrelevant.

## Decisions

### Refuse-on-dirty over auto-stash

If any worktree has uncommitted changes (tracked or untracked), `wrt rebase` lists them and exits non-zero without touching anything. The user stashes/commits and re-runs.

**Why**: auto-stash + auto-pop is convenient until a conflict during pop leaves the user in a confusing intermediate state. Refusing is honest; the user already knows whether their work-in-progress is worth saving.

**Implementation**: per worktree, run `git status --porcelain`. Empty output = clean. Non-empty = dirty, list it.

### Use `<base>@{upstream}` for fetch + rebase target

For each worktree, look up its base branch (`RepoState.BaseBranch` for code, `BackportState.BaseBranch` for backports). Resolve its upstream:
```
git -C <base-repo> rev-parse --abbrev-ref <base>@{upstream}
```
This returns something like `origin/main` or `upstream/release/1.1.1`, or fails with exit code if no upstream is set. Split on `/` to get the remote name and the ref.

Then in the base repo:
```
git -C <base-repo> fetch <remote>
```
Then in the worktree:
```
git -C <worktree> rebase <remote>/<branch>
```

**No upstream → skip that worktree, continue**. The skip is reported but does not halt the run, because a missing upstream is a configuration issue specific to that branch, not a failure of the rebase machinery.

**Why not `git pull --rebase`**: pull does fetch + rebase but couples the two; and it has historically had subtle behaviors with non-default remotes. Explicit fetch + rebase is clearer and easier to report on.

### Stop-on-conflict for actual conflicts

If `git rebase` exits non-zero (and it's not the "no upstream" case handled above), assume it's a conflict (or some other rebase-time failure). Leave the worktree in whatever state git left it — typically mid-rebase with conflict markers. Print:
- which worktree
- which files are conflicted (parse `git status --porcelain` or `git diff --name-only --diff-filter=U`)
- the resolution path: `cd <worktree>` → resolve → `git rebase --continue` → re-run `wrt rebase`

Then bail. Don't try the next worktree — the user's attention is already on this one.

**Why stop**: continuing would risk piling up multiple mid-rebase states across worktrees, which is hard to reason about. Better one problem at a time.

### Current task only

The command finds the current task via `taskctx.Find` (same as `wrt repo add`, `wrt backport add`). No `wrt rebase <task-id>` and no `--all`. If multi-task or cross-task scope becomes useful, add it later — v1 stays small.

### Summary output

End of run, print a single summary block:
```
Rebase summary for task IR-42:
  some-repo/code (users/db/IR-42 ← origin/main)              ok, 3 commits replayed
  some-repo/backports/1.1.1 (← origin/release/1.1.1)         ok, no changes
  other-repo/code (users/db/IR-42 ← upstream/main)           CONFLICT
    files: src/foo.go, src/bar.go
    cd repositories/other-repo/code && resolve, then git rebase --continue
  third-repo/code (users/db/IR-42 ← ?)                       SKIPPED — no upstream tracking
```

If the run halted on conflict, the summary covers worktrees attempted; remaining worktrees are listed as `pending`.

If the run was refused on dirty worktrees, the output is a different shape — just the dirty list and a hint.

### Command name: `rebase`

The user picked `rebase` over `refresh` because the operation is unambiguously a rebase. If we ever add a fetch-only mode or a merge mode, a sibling command (`wrt fetch`, `wrt sync`) or a flag (`wrt rebase --merge`) is the natural extension.

## Risks / Trade-offs

- **Stops on first conflict.** Multi-repo tasks where conflicts are isolated to one repo don't get full multi-repo progress in one shot. Trade-off accepted; the user runs `wrt rebase` again after resolving.
- **No upstream skip.** A task branch with a base on a local-only branch (no tracking remote) is silently un-rebased. We report it in the summary, but a user who doesn't read the summary might assume they're up to date. Mitigation: make the skip-line in the summary visually distinct (e.g. SKIPPED in caps or with a bullet).
- **Fetch failures.** `git fetch` can fail for network reasons. Treat as a hard error for that worktree, same handling as a rebase conflict (stop, report, bail).
- **Conflicts during fetch resolution.** Not really a thing for fetch; reservations apply to rebase only.
- **Sub-second runtimes feel like nothing happened.** Each worktree should print a line as it goes, not just at the summary, so the user sees progress.
