## MODIFIED Requirements

### Requirement: task.json schema
Each `task.json` SHALL contain:
- `name` (string, required): task identifier (e.g. `CLOUD-111`)
- `url` (string, optional): link to issue tracker
- `description` (string, optional): short description
- `created` (string, required): ISO 8601 timestamp of task creation
- `repositories` (object, optional): map of repo name to repo state:
  - `base_branch`: branch the task branch was created from
  - `task_branch`: full branch name
  - `backport_branches`: list of version aliases with active backport worktrees

The repo state SHALL NOT include the base repository's on-disk path. The path is resolved at command time by looking up the repo in the registry (`<task-root>/repos/<name>/repo.json`), which is the single source of truth.

#### Scenario: Minimal task.json
- **WHEN** a task is created with only a name
- **THEN** `task.json` contains `name`, `created`, and empty `repositories` object

#### Scenario: Repo state contains no base repo path
- **WHEN** `wrt repo add some-repo` writes the repo state
- **THEN** the entry for `some-repo` contains only `base_branch`, `task_branch`, and `backport_branches` — no `repo_path` field

### Requirement: Close task
`wrt close <task-name>` SHALL delete all git worktrees associated with the task, generate patch files for each repo, move the task directory from `tasks/open/` to `tasks/closed/`, and remove the now-empty `code/` and `backports/` directories. Closing SHALL NOT require the task branch to be merged. For each repo in `task.json`, the base repository path used to invoke `git worktree remove` SHALL be resolved from the registry by repo name.

#### Scenario: Close with repos
- **WHEN** `wrt close JIRA-111` is run and the task has two repos
- **THEN** patch files are written, worktrees are removed via `git worktree remove`, and the task moves to `tasks/closed/`

#### Scenario: Close already closed task
- **WHEN** `wrt close` is run on a task in `tasks/closed/`
- **THEN** the tool errors: "task 'JIRA-111' is already closed"

#### Scenario: Repo missing from registry
- **WHEN** `wrt close JIRA-111` is run and a repo listed in `task.json` is no longer in the registry
- **THEN** the tool errors with a clear message identifying which repo cannot be resolved

### Requirement: Reopen task
`wrt reopen <task-name>` SHALL move a closed task from `tasks/closed/` to `tasks/open/`, recreate git worktrees for each repo on the task branch, and attempt to apply saved patch files via `git am`. On patch conflict, the tool SHALL leave the worktree in `git am` paused state and print instructions for manual resolution without aborting. For each repo, the base repository path used to invoke `git worktree add` SHALL be resolved from the registry by repo name.

#### Scenario: Clean reopen
- **WHEN** patches apply cleanly
- **THEN** the task is fully restored to open state with all worktrees functional

#### Scenario: Patch conflict on reopen
- **WHEN** `git am` encounters conflicts
- **THEN** the tool prints the conflicting files, instructions to resolve and run `git am --continue`, and exits non-zero without running `git am --abort`

#### Scenario: Reopen open task
- **WHEN** `wrt reopen` is run on a task already in `tasks/open/`
- **THEN** the tool errors: "task 'JIRA-111' is already open"

#### Scenario: Repo missing from registry
- **WHEN** `wrt reopen JIRA-111` is run and a repo listed in `task.json` is no longer in the registry
- **THEN** the tool errors with a clear message identifying which repo cannot be resolved
