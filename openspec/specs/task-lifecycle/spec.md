## ADDED Requirements

### Requirement: Task directory structure
An open task SHALL be a directory at `<task-root>/tasks/open/<task-id>/` containing `task.json`, `AGENTS.md`, `ac.md`, and a `repositories/` subdirectory. A closed task SHALL live at `<task-root>/tasks/closed/<task-id>/` with the same files but an empty or patch-only `repositories/` directory.

#### Scenario: Open task layout
- **WHEN** a task is created and repos are added
- **THEN** the directory structure matches `tasks/open/<task-id>/repositories/<repo-name>/code/`

#### Scenario: Closed task layout
- **WHEN** a task is closed
- **THEN** it moves to `tasks/closed/<task-id>/` with worktrees removed

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

#### Scenario: Minimal task.json
- **WHEN** a task is created with only a name
- **THEN** `task.json` contains `name`, `created`, and empty `repositories` object

### Requirement: Create task
`wrt create` SHALL create a new task workspace. It SHALL support interactive prompts and CLI flags. The task name is required; URL and description are optional. It SHALL error if a task with that name already exists in open or closed. The task name SHALL NOT contain spaces or forward slashes; `wrt create` SHALL reject such names with a clear error before any filesystem operations.

#### Scenario: Interactive create
- **WHEN** `wrt create` is run with no arguments
- **THEN** the tool prompts for name (required), URL (optional), description (optional)

#### Scenario: Name conflict
- **WHEN** `wrt create CLOUD-111` is run and `CLOUD-111` already exists (open or closed)
- **THEN** the tool errors with the task's current location

#### Scenario: Task created
- **WHEN** a task is successfully created
- **THEN** `task.json`, `AGENTS.md`, and `ac.md` are written, and the tool prints the task path

#### Scenario: Name contains space
- **WHEN** `wrt create "my task"` is run
- **THEN** the tool errors: task name must not contain spaces or slashes

#### Scenario: Name contains slash
- **WHEN** `wrt create "JIRA/111"` is run
- **THEN** the tool errors: task name must not contain spaces or slashes

### Requirement: List tasks
`wrt list` SHALL list all open tasks ordered by `created` timestamp descending (newest first), showing: name, description (if set), created date, and number of repos attached.

#### Scenario: No open tasks
- **WHEN** `tasks/open/` is empty
- **THEN** the tool prints "No open tasks."

#### Scenario: Multiple tasks
- **WHEN** multiple open tasks exist
- **THEN** they are displayed newest-first with name, created date, and repo count

### Requirement: Task info
`wrt info <task-name>` SHALL display full task details: name, URL, description, created date, and for each repo: task branch, base branch, and list of backport versions with their branch names.

#### Scenario: Task not found
- **WHEN** the specified task name does not exist in open or closed
- **THEN** the tool errors: "task '<name>' not found"

#### Scenario: Worktree discrepancy
- **WHEN** a repo is listed in `task.json` but its worktree directory does not exist on disk
- **THEN** `wrt info` flags it as "⚠ worktree missing"

### Requirement: Task path
`wrt path <arg>` SHALL print only the absolute path to the task directory, with no trailing newline decoration, suitable for use in `cd $(wrt path <arg>)`. The argument SHALL be treated as either a bare task name or a formatted `wrt list` output line; in either case the task name is the first whitespace-delimited word of the argument.

#### Scenario: Valid task by name
- **WHEN** `wrt path task-foo` is run and the task exists
- **THEN** only the absolute path is printed to stdout, no other output

#### Scenario: Valid task from list line
- **WHEN** `wrt path "task-foo   2025-01-10   2   fix the login bug"` is run and task-foo exists
- **THEN** only the absolute path for task-foo is printed to stdout

#### Scenario: Task not found
- **WHEN** the first word of the argument does not match any task
- **THEN** the tool errors to stderr and exits non-zero

#### Scenario: Empty argument
- **WHEN** `wrt path ""` is run
- **THEN** the tool errors to stderr and exits non-zero

### Requirement: Close task
`wrt close <task-name>` SHALL delete all git worktrees associated with the task, generate patch files for each repo, move the task directory from `tasks/open/` to `tasks/closed/`, and remove the now-empty `code/` and `backports/` directories. Closing SHALL NOT require the task branch to be merged.

#### Scenario: Close with repos
- **WHEN** `wrt close JIRA-111` is run and the task has two repos
- **THEN** patch files are written, worktrees are removed via `git worktree remove`, and the task moves to `tasks/closed/`

#### Scenario: Close already closed task
- **WHEN** `wrt close` is run on a task in `tasks/closed/`
- **THEN** the tool errors: "task 'JIRA-111' is already closed"

### Requirement: Reopen task
`wrt reopen <task-name>` SHALL move a closed task from `tasks/closed/` to `tasks/open/`, recreate git worktrees for each repo on the task branch, and attempt to apply saved patch files via `git am`. On patch conflict, the tool SHALL leave the worktree in `git am` paused state and print instructions for manual resolution without aborting.

#### Scenario: Clean reopen
- **WHEN** patches apply cleanly
- **THEN** the task is fully restored to open state with all worktrees functional

#### Scenario: Patch conflict on reopen
- **WHEN** `git am` encounters conflicts
- **THEN** the tool prints the conflicting files, instructions to resolve and run `git am --continue`, and exits non-zero without running `git am --abort`

#### Scenario: Reopen open task
- **WHEN** `wrt reopen` is run on a task already in `tasks/open/`
- **THEN** the tool errors: "task 'JIRA-111' is already open"
