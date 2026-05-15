## ADDED Requirements

### Requirement: AGENTS.md auto-generated at task creation
`wrt create` SHALL generate an `AGENTS.md` file in the task root describing the workspace structure and agent conventions. The content SHALL be static (not user-editable at creation time) and SHALL explain folder layout, the purpose of each context file, and the log.md convention.

#### Scenario: AGENTS.md created
- **WHEN** `wrt create JIRA-111` succeeds
- **THEN** `tasks/open/JIRA-111/AGENTS.md` exists with generated content

#### Scenario: AGENTS.md content — structure section
- **WHEN** an agent reads AGENTS.md
- **THEN** it finds a section describing the folder layout: `repositories/<repo>/code/` for main worktrees and `repositories/<repo>/backports/<version>/` for backport worktrees

#### Scenario: AGENTS.md content — context files section
- **WHEN** an agent reads AGENTS.md
- **THEN** it finds descriptions of: `ac.md` (acceptance criteria), `log.md` (agent activity log)

#### Scenario: AGENTS.md content — log.md instruction
- **WHEN** an agent reads AGENTS.md
- **THEN** it finds an explicit instruction to record all significant actions in `log.md` with a timestamp and short description of what was done and why

### Requirement: ac.md scaffolded at task creation
`wrt create` SHALL create an `ac.md` file in the task root with a minimal template (header and placeholder). The file is intended for the user to fill in acceptance criteria.

#### Scenario: ac.md created
- **WHEN** `wrt create JIRA-111` succeeds
- **THEN** `tasks/open/JIRA-111/ac.md` exists with a placeholder template

### Requirement: Context files preserved on close and reopen
`AGENTS.md` and `ac.md` SHALL be preserved when a task is closed and SHALL remain present when a task is reopened.

#### Scenario: Files survive close
- **WHEN** `wrt close JIRA-111` runs
- **THEN** `AGENTS.md` and `ac.md` are present in `tasks/closed/JIRA-111/`

#### Scenario: Files survive reopen
- **WHEN** `wrt reopen JIRA-111` runs
- **THEN** `AGENTS.md` and `ac.md` are present in `tasks/open/JIRA-111/`
