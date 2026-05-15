## REMOVED Requirements

### Requirement: AGENTS.md auto-generated at task creation
**Reason**: AGENTS.md is no longer generated from a hardcoded Go string literal. It is now sourced from `<task-root>/tasks/task-template/AGENTS.md`, copied with variable substitution by the template mechanism (see capability `task-templates`).
**Migration**: The default content shipped in the template (written by `wrt init` or lazily on first `wrt create`) is equivalent to the previous generator output. Existing tasks are not retroactively updated.

### Requirement: ac.md scaffolded at task creation
**Reason**: ac.md is no longer scaffolded from a hardcoded literal. It is now sourced from `<task-root>/tasks/task-template/ac.md`.
**Migration**: Same as AGENTS.md — default content is preserved via the prepopulated template.

## MODIFIED Requirements

### Requirement: Context files preserved on close and reopen
`AGENTS.md`, `ac.md`, `.claude/settings.local.json`, and any other files originally scaffolded from the task template SHALL be preserved when a task is closed and SHALL remain present when a task is reopened.

#### Scenario: Files survive close
- **WHEN** `wrt close JIRA-111` runs
- **THEN** `AGENTS.md`, `ac.md`, and `.claude/settings.local.json` are present in `tasks/closed/JIRA-111/`

#### Scenario: Files survive reopen
- **WHEN** `wrt reopen JIRA-111` runs
- **THEN** all task-template-sourced files are present in `tasks/open/JIRA-111/`
