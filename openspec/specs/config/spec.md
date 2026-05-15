## ADDED Requirements

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
