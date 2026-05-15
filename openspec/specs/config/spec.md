## Purpose
Provides the tool's global configuration surface: task root location, optional username override, and discovery of `~/.config/wrt/config.toml`.
## Requirements
### Requirement: Global configuration file
The tool SHALL read its configuration from `~/.config/wrt/config.toml`. If the file does not exist, the tool SHALL use built-in defaults and SHALL NOT error.

#### Scenario: Config file missing
- **WHEN** `config.toml` does not exist
- **THEN** the tool uses default values for all settings and runs normally

#### Scenario: Config file present
- **WHEN** `config.toml` exists and contains a `task_root` key
- **THEN** the tool uses that path as the task root for all operations

### Requirement: Task root configuration
The tool SHALL support a configurable `task_root` path that defines where all task workspaces and the repo registry are stored. There is no default — if `task_root` is not configured and no task root exists, the tool SHALL prompt the user to run `wrt init` or set the path.

#### Scenario: Task root not configured
- **WHEN** `task_root` is not set in config and the tool is invoked
- **THEN** the tool prints an actionable error message explaining how to set it

#### Scenario: Task root configured
- **WHEN** `task_root` is set to a valid directory path
- **THEN** all commands resolve task and repo paths relative to that root

### Requirement: Username configuration
The tool SHALL support an optional `username` field in `config.toml` that overrides `$USER` when substituting `{user}` in branch name templates. If not set, `$USER` SHALL be used.

#### Scenario: Username not configured
- **WHEN** `username` is absent from config
- **THEN** `$USER` environment variable is used for `{user}` substitution

#### Scenario: Username configured
- **WHEN** `username` is set in config
- **THEN** that value is used for `{user}` substitution in all branch templates

### Requirement: Username sanitization
When substituting `{user}` into branch name templates, the tool SHALL strip whitespace and lowercase the username to ensure valid git branch names.

#### Scenario: Username with spaces
- **WHEN** the resolved username contains spaces (e.g. "Daniel Barton")
- **THEN** the substituted value is "danielbarton" (stripped and lowercased)

### Requirement: Init command
The tool SHALL provide a `wrt init` command that prepares the configured task root for use. The command SHALL be idempotent: running it on an already-initialized task root SHALL produce no destructive changes.

`wrt init` SHALL:
- Resolve the task root from configuration; error with an actionable message if `task_root` is unset.
- Create `tasks/open/`, `tasks/closed/`, `repos/`, and `tasks/task-template/` under the task root if any are missing.
- Populate `tasks/task-template/` with default content (AGENTS.md, ac.md, `.claude/settings.local.json`) for files that do not yet exist.
- Print a short summary of what was created versus what was already present.

#### Scenario: Init on fresh task root
- **WHEN** `wrt init` is run and the task root contains none of the expected subdirectories
- **THEN** `tasks/open/`, `tasks/closed/`, `repos/`, and `tasks/task-template/` are created, default template files are written, and the summary lists them as created

#### Scenario: Init re-run
- **WHEN** `wrt init` is run on a fully initialized task root
- **THEN** no files are modified; the summary reports everything as already present

#### Scenario: Init without task root
- **WHEN** `wrt init` is run and `task_root` is not configured
- **THEN** the tool errors with the same actionable message documented under "Task root configuration"

#### Scenario: Init preserves custom template
- **WHEN** `wrt init` is run and the user has previously edited `tasks/task-template/AGENTS.md`
- **THEN** the edited file is not modified; any other missing default template files are still written

