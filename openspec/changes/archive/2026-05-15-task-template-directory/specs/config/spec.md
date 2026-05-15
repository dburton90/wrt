## ADDED Requirements

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
