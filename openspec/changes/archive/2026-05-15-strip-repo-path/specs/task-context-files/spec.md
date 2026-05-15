## MODIFIED Requirements

### Requirement: AGENTS.md auto-generated at task creation
`wrt create` SHALL generate an `AGENTS.md` file in the task root describing the workspace structure and agent conventions. The content SHALL be static (not user-editable at creation time) and SHALL explain folder layout, the purpose of each context file, the log.md convention, and explicitly state that work happens inside the per-repo worktrees and SHALL NOT touch the base repository on disk.

#### Scenario: AGENTS.md created
- **WHEN** `wrt create JIRA-111` succeeds
- **THEN** `tasks/open/JIRA-111/AGENTS.md` exists with generated content

#### Scenario: AGENTS.md content — structure section
- **WHEN** an agent reads AGENTS.md
- **THEN** it finds a section describing the folder layout: `repositories/<repo>/code/` for main worktrees and `repositories/<repo>/backports/<version>/` for backport worktrees

#### Scenario: AGENTS.md content — workspace boundary
- **WHEN** an agent reads AGENTS.md
- **THEN** it finds an explicit statement that all work happens inside `repositories/<repo>/code/` (or the relevant backport worktree) and that the base repository registered in the wrt registry is internal plumbing the agent SHALL NOT modify directly

#### Scenario: AGENTS.md content — context files section
- **WHEN** an agent reads AGENTS.md
- **THEN** it finds descriptions of: `ac.md` (acceptance criteria), `log.md` (agent activity log)

#### Scenario: AGENTS.md content — log.md instruction
- **WHEN** an agent reads AGENTS.md
- **THEN** it finds an explicit instruction to record all significant actions in `log.md` with a timestamp and short description of what was done and why
