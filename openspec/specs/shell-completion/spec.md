## Purpose
Provides shell tab-completion scripts and the `wrt-cd` fzf-based task navigator function.

## Requirements

### Requirement: Completion generation command
The tool SHALL provide a `wrt completion <shell>` command that outputs a completion script for the specified shell. Supported shells SHALL be: `bash`, `zsh`, `fish`. For `bash` and `zsh`, the output SHALL include a `wrt-cd` shell function definition appended after the cobra-generated completion script. The function SHALL pipe `wrt list` output to `fzf --header-lines=2`, and on a non-empty selection SHALL `cd` to the path returned by `wrt path`. Fish output SHALL NOT include `wrt-cd` (deferred).

#### Scenario: Generate zsh completion
- **WHEN** `wrt completion zsh` is run
- **THEN** a valid zsh completion script is printed to stdout, followed by the `wrt-cd` function definition

#### Scenario: Generate bash completion
- **WHEN** `wrt completion bash` is run
- **THEN** a valid bash completion script is printed to stdout, followed by the `wrt-cd` function definition

#### Scenario: wrt-cd navigates to selected task
- **WHEN** the user runs `wrt-cd` and selects a task from the fzf picker
- **THEN** the shell changes directory to that task's directory

#### Scenario: wrt-cd cancelled
- **WHEN** the user runs `wrt-cd` and cancels fzf (ctrl-c or no selection)
- **THEN** the current directory is unchanged and no error is printed

### Requirement: Task name completion
Commands that accept `<task-name>` as an argument SHALL complete to the list of task names from `tasks/open/` (and `tasks/closed/` where relevant).

#### Scenario: wrt info completion
- **WHEN** user types `wrt info <TAB>`
- **THEN** completion suggests all open task names

#### Scenario: wrt close completion
- **WHEN** user types `wrt close <TAB>`
- **THEN** completion suggests all open task names

#### Scenario: wrt reopen completion
- **WHEN** user types `wrt reopen <TAB>`
- **THEN** completion suggests all closed task names

### Requirement: Repo name completion
Commands that accept `<repo-name>` SHALL complete to the list of registered repos from `<task-root>/repos/`.

#### Scenario: wrt repo add completion
- **WHEN** user types `wrt repo add <TAB>`
- **THEN** completion suggests all registered repo names

#### Scenario: wrt backport add repo completion
- **WHEN** user types `wrt backport add <TAB>`
- **THEN** completion suggests repo names present in the current task's `task.json`

### Requirement: Backport version completion
The second argument of `wrt backport add <repo> <version>` SHALL complete to the list of backport versions configured for the specified repo in `repo.json`.

#### Scenario: wrt backport add version completion
- **WHEN** user types `wrt backport add some-repo <TAB>`
- **THEN** completion suggests the version aliases from `some-repo/repo.json` backport_branches keys (e.g. `1.1.1`, `6.9.x`)
