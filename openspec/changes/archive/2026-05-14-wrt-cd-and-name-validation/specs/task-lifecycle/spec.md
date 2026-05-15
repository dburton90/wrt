## MODIFIED Requirements

### Requirement: Create task
`wrt create` SHALL create a new task workspace. It SHALL support interactive prompts and CLI flags. The task name is required; URL and description are optional. It SHALL error if a task with that name already exists in open or closed. The task name SHALL NOT contain spaces or forward slashes; `wrt create` SHALL reject such names with a clear error before any filesystem operations.

#### Scenario: Interactive create
- **WHEN** `wrt create` is run with no arguments
- **THEN** the tool prompts for name (required), URL (optional), description (optional)

#### Scenario: Name conflict
- **WHEN** `wrt create JIRA-111` is run and `JIRA-111` already exists (open or closed)
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
