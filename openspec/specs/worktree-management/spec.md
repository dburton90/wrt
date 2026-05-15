## Purpose
Adds repository worktrees to a task on a freshly created task branch, detects the current task by walking up from CWD, and pins worktree paths to `repositories/<repo>/code/`.

## Requirements

### Requirement: Task context detection
Commands that operate on a task (`wrt repo add`, `wrt backport add`) SHALL detect the current task by traversing up from CWD looking for `task.json`, stopping at the filesystem root. If no `task.json` is found, the command SHALL error with a clear message.

#### Scenario: Run from inside task
- **WHEN** CWD is `tasks/open/JIRA-111/repositories/some-repo/code/src/`
- **THEN** the tool finds `task.json` at `tasks/open/JIRA-111/task.json` and operates on task `JIRA-111`

#### Scenario: Run outside any task
- **WHEN** no `task.json` exists in CWD or any ancestor
- **THEN** the tool errors: "not inside a task directory. Navigate into a task or use `wrt list` to find one."

### Requirement: Add repo worktree to task
`wrt repo add <repo-name>` SHALL look up the repo in the registry, create a new git branch from `default_base_branch` using the task branch template, create a git worktree at `repositories/<repo-name>/code/`, and update `task.json` with the repo entry.

#### Scenario: Repo not in registry
- **WHEN** `wrt repo add unknown-repo` is run
- **THEN** the tool errors: "repo 'unknown-repo' not found. Run `wrt repo list` to see registered repos or `wrt repo create unknown-repo` to add it."

#### Scenario: Repo already added to task
- **WHEN** `wrt repo add some-repo` is run and some-repo is already in `task.json`
- **THEN** the tool errors: "repo 'some-repo' is already added to this task"

#### Scenario: Successful repo add
- **WHEN** `wrt repo add some-repo` is run inside a valid task
- **THEN** the branch `users/dbarton/JIRA-111` is created from `main`, a worktree is created at `repositories/some-repo/code/`, and `task.json` is updated

### Requirement: Worktree at known path
The git worktree for a repo SHALL always be created at `<task-dir>/repositories/<repo-name>/code/` regardless of the repo's own directory name or location.

#### Scenario: Worktree path
- **WHEN** repo `some-repo` is added to task `JIRA-111`
- **THEN** the worktree exists at `tasks/open/JIRA-111/repositories/some-repo/code/`

### Requirement: task.json updated on repo add
After successfully creating the worktree, `task.json` SHALL be updated to record: `base_branch`, `task_branch` (resolved from template), and an empty `backport_branches` list.

#### Scenario: task.json after add
- **WHEN** `wrt repo add some-repo` succeeds
- **THEN** `task.json` contains an entry for `some-repo` with the resolved branch name
