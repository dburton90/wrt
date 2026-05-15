# worktree-rebase Specification

## Purpose
TBD - created by archiving change wrt-rebase. Update Purpose after archive.
## Requirements
### Requirement: Rebase current task
The tool SHALL provide a `wrt rebase` command that operates on the current task (detected via the existing task-context lookup). The command SHALL iterate every worktree of the task — the `code/` worktree for each repo and each `backports/<version>/` worktree — and rebase each onto its base branch's upstream tracking ref.

#### Scenario: Outside a task
- **WHEN** `wrt rebase` is run and no `task.json` is found in CWD or any ancestor
- **THEN** the tool errors with the same message used by other task-context commands

#### Scenario: Successful rebase
- **WHEN** every worktree is clean, has a base branch with an upstream, and rebases without conflict
- **THEN** every worktree advances to the upstream ref and the tool prints an ok line for each

### Requirement: Refuse on dirty worktrees
Before performing any rebase, the tool SHALL check every worktree associated with the task using `git status --porcelain`. If any worktree reports uncommitted changes (tracked or untracked), the tool SHALL print the list of dirty worktrees and exit non-zero without performing fetches or rebases.

#### Scenario: One dirty worktree
- **WHEN** `wrt rebase` is run and exactly one of the task's worktrees has uncommitted changes
- **THEN** the tool lists that worktree, performs no fetch or rebase on any worktree, and exits non-zero

#### Scenario: All clean
- **WHEN** every worktree is clean
- **THEN** the pre-flight check passes silently and the rebase loop proceeds

### Requirement: Upstream resolution
For each worktree, the tool SHALL resolve the upstream of its base branch using git's upstream tracking (`<base>@{upstream}` via `git rev-parse --abbrev-ref`). The tool SHALL NOT assume any particular remote name (e.g. `origin`). It SHALL split the returned `<remote>/<ref>` into the remote name and the ref name, and SHALL fetch that remote in the base repository before rebasing onto `<remote>/<ref>`.

#### Scenario: Origin remote
- **WHEN** the base branch's upstream is `origin/main`
- **THEN** the tool runs `git fetch origin` and rebases the worktree onto `origin/main`

#### Scenario: Non-origin remote
- **WHEN** the base branch's upstream is `upstream/release/1.1.1`
- **THEN** the tool runs `git fetch upstream` and rebases the worktree onto `upstream/release/1.1.1`

#### Scenario: No upstream
- **WHEN** the base branch has no upstream tracking ref
- **THEN** the tool reports that worktree as "no upstream tracking", skips it, and continues with the remaining worktrees; the run as a whole does not fail

### Requirement: Stop on conflict
When a `git rebase` invocation exits non-zero (other than the no-upstream case handled above), the tool SHALL leave the worktree in whatever state `git rebase` left it, print the worktree path and the list of files with merge conflicts, print explicit instructions to `cd` into the worktree, resolve the conflicts, run `git rebase --continue`, and re-run `wrt rebase` to continue with remaining worktrees. The tool SHALL NOT attempt subsequent worktrees in the same invocation and SHALL exit non-zero.

#### Scenario: Conflict on first worktree
- **WHEN** the first worktree's rebase conflicts
- **THEN** the tool reports the conflict and the remaining worktrees are listed as pending in the summary; the tool exits non-zero

#### Scenario: Conflict on a later worktree
- **WHEN** a conflict occurs after one or more earlier worktrees have rebased successfully
- **THEN** the successful ones are reported as ok and the remaining (including those after the conflict) are reported as pending

#### Scenario: Re-run after resolution
- **WHEN** the user resolves the conflict, runs `git rebase --continue`, and re-runs `wrt rebase`
- **THEN** the previously-conflicted worktree is now clean (already rebased) and the remaining pending worktrees are processed

### Requirement: Summary output
At the end of every run (whether successful, halted on conflict, or refused on dirty), the tool SHALL print a summary identifying each worktree by repo and role (code / backport-<version>), the base ref it was rebased onto (or would have been), and the outcome: `ok`, `CONFLICT`, `SKIPPED — no upstream`, `SKIPPED — dirty`, or `pending`.

#### Scenario: Mixed summary
- **WHEN** a run includes successes, one conflict, and one no-upstream skip
- **THEN** the summary lists every target with the appropriate status keyword

