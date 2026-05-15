## MODIFIED Requirements

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
