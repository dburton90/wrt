# task-templates Specification

## Purpose
TBD - created by archiving change task-template-directory. Update Purpose after archive.
## Requirements
### Requirement: Template directory location
The task root SHALL contain a template directory at `<task-root>/tasks/task-template/`. Files and subdirectories under this path SHALL be treated as the scaffolding source for new tasks. The directory SHALL be created and prepopulated either by `wrt init` or lazily on the first `wrt create` that finds it missing.

#### Scenario: Template directory created by init
- **WHEN** `wrt init` is run on a task root that has no `tasks/task-template/`
- **THEN** the directory is created and populated with the default template files

#### Scenario: Template directory created lazily
- **WHEN** `wrt create JIRA-111` is run and `tasks/task-template/` does not exist
- **THEN** the directory is created and populated with defaults before the task is scaffolded

#### Scenario: Template directory ignored by task listing
- **WHEN** `wrt list` is run
- **THEN** `tasks/task-template/` is not listed as a task (it has no `task.json`)

### Requirement: Default template content
When prepopulating `tasks/task-template/`, the tool SHALL write the following files only if they do not already exist:
- `AGENTS.md` — agent workspace orientation, including folder layout and the worktree-vs-base-repo boundary
- `ac.md` — acceptance criteria placeholder
- `.claude/settings.local.json` — minimal Claude Code permissions allowlist with relative and absolute `Bash(cd …)` entries

#### Scenario: All defaults written on first init
- **WHEN** `wrt init` runs on an empty `tasks/task-template/`
- **THEN** `AGENTS.md`, `ac.md`, and `.claude/settings.local.json` are written with default content

#### Scenario: Existing files preserved
- **WHEN** `wrt init` runs and `tasks/task-template/AGENTS.md` already exists with custom content
- **THEN** the existing file is not modified; other missing defaults are still written

### Requirement: Template applied on task creation
`wrt create` SHALL apply the template by walking `tasks/task-template/`, mirroring its directory structure into the new task directory, copying each file with variable substitution applied to its contents. Destination files that already exist SHALL NOT be overwritten.

#### Scenario: Template files copied
- **WHEN** `wrt create JIRA-111` succeeds and `tasks/task-template/` contains `AGENTS.md`, `ac.md`, and `.claude/settings.local.json`
- **THEN** `tasks/open/JIRA-111/` contains those same files with substituted content

#### Scenario: User-added template files included
- **WHEN** the user has added `Makefile` to `tasks/task-template/`
- **THEN** `wrt create` copies `Makefile` into every new task

#### Scenario: Skip-if-exists
- **WHEN** a file already exists in the destination task directory
- **THEN** the template copy does not overwrite it

### Requirement: Variable substitution
The tool SHALL substitute the following variables in template file *contents* (not filenames or directory names):
- `{task-id}` — the task name (e.g. `JIRA-111`)
- `{task-dir}` — absolute path to the new task directory
- `{task-root}` — absolute path to the task root
- `{user}` — resolved username (same sanitization as branch-template `{user}`)

Substitution is plain text replacement. The tool SHALL NOT interpret any other syntax.

#### Scenario: task-id substituted
- **WHEN** a template file contains `# Task Workspace: {task-id}` and the task is `JIRA-111`
- **THEN** the destination file contains `# Task Workspace: JIRA-111`

#### Scenario: task-dir substituted in allowlist
- **WHEN** `.claude/settings.local.json` contains `Bash(cd {task-dir}/repositories/*)`
- **THEN** the destination file contains the absolute path to the task directory

#### Scenario: Filenames not substituted
- **WHEN** a template directory contains a file literally named `{task-id}.md`
- **THEN** the destination contains a file named `{task-id}.md` (no rename)

#### Scenario: Unknown variable left intact
- **WHEN** a template contains `{unknown-var}`
- **THEN** the destination contains `{unknown-var}` literally

