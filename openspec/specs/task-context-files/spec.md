## Purpose
Ensures every task workspace ships with agent-facing context files (AGENTS.md, ac.md, etc.) that describe layout and conventions.
## Requirements
### Requirement: Context files preserved on close and reopen
`AGENTS.md`, `ac.md`, `.claude/settings.local.json`, and any other files originally scaffolded from the task template SHALL be preserved when a task is closed and SHALL remain present when a task is reopened.

#### Scenario: Files survive close
- **WHEN** `wrt close JIRA-111` runs
- **THEN** `AGENTS.md`, `ac.md`, and `.claude/settings.local.json` are present in `tasks/closed/JIRA-111/`

#### Scenario: Files survive reopen
- **WHEN** `wrt reopen JIRA-111` runs
- **THEN** all task-template-sourced files are present in `tasks/open/JIRA-111/`

