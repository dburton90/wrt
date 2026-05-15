## MODIFIED Requirements

### Requirement: Create task
`wrt create` SHALL create a new task workspace. It SHALL support interactive prompts and CLI flags. The task name is required; URL and description are optional. It SHALL error if a task with that name already exists in open or closed. The task name SHALL NOT contain spaces or forward slashes; `wrt create` SHALL reject such names with a clear error before any filesystem operations.

After the task directory is created and `task.json` written, the tool SHALL ensure `tasks/task-template/` exists (initializing it with default content if missing), then scaffold the new task by applying the template per the `task-templates` capability — copying every file from the template directory into the new task directory with variable substitution.

#### Scenario: Interactive create
- **WHEN** `wrt create` is run with no arguments
- **THEN** the tool prompts for name (required), URL (optional), description (optional)

#### Scenario: Name conflict
- **WHEN** `wrt create CLOUD-111` is run and `CLOUD-111` already exists (open or closed)
- **THEN** the tool errors with the task's current location

#### Scenario: Task created
- **WHEN** a task is successfully created
- **THEN** `task.json` is written, every file in `tasks/task-template/` is copied into the new task directory with substitution, and the tool prints the task path

#### Scenario: Name contains space
- **WHEN** `wrt create "my task"` is run
- **THEN** the tool errors: task name must not contain spaces or slashes

#### Scenario: Name contains slash
- **WHEN** `wrt create "JIRA/111"` is run
- **THEN** the tool errors: task name must not contain spaces or slashes

#### Scenario: Template directory missing
- **WHEN** `wrt create JIRA-111` is run and `tasks/task-template/` does not yet exist
- **THEN** the template directory is created and prepopulated with defaults before scaffolding proceeds
