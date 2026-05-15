## Why

`AGENTS.md` and `ac.md` are currently hardcoded Go string literals inside `cmd/contextfiles.go`. Adding any new per-task scaffolded file (e.g. `.claude/settings.local.json` to seed agent permissions, or an `.editorconfig`, or a Makefile snippet) means a code change and a rebuild. Worse, the content is opinionated — what one task root's agent setup needs is not what another's needs — but with hardcoded content there's no per-user knob.

A file-based template directory at `<task-root>/tasks/task-template/` solves both: anything dropped in there is copied into each new task, with simple variable substitution for paths and names. The user owns the content; `wrt` owns the placement.

## What Changes

- New directory `<task-root>/tasks/task-template/`.
- New command `wrt init`: creates the task-root skeleton (`tasks/open/`, `tasks/closed/`, `repos/`, `tasks/task-template/`) and prepopulates `task-template/` with default content (AGENTS.md, ac.md, .claude/settings.local.json).
- `wrt create` lazily creates and prepopulates `task-template/` if it does not yet exist, then proceeds.
- `wrt create` copies every file under `task-template/` into the new task directory, applying substitution of `{task-id}`, `{task-dir}`, `{task-root}`, `{user}` in file contents (not file names). Existing files in the destination are left untouched (skip-if-exists).
- `cmd/contextfiles.go` no longer hardcodes AGENTS.md/ac.md content — that content moves into the default template baked into the binary as the prepopulated content.
- Existing tasks are NOT retroactively scaffolded; only new tasks see the templates.

## Capabilities

### New Capabilities

- `task-templates`: file-based scaffolding mechanism — `tasks/task-template/` directory and substitution rules; how `wrt create` applies it.

### Modified Capabilities

- `config`: introduces `wrt init` for explicit task-root setup.
- `task-lifecycle`: `wrt create` populates the task directory by copying from `task-template/` (with substitution) instead of writing hardcoded files.
- `task-context-files`: AGENTS.md and ac.md remain part of every task, but they originate from the template directory rather than from baked-in string literals. Behaviour is preserved; mechanism changes.

## Impact

- New package `internal/template` (or similar) for copy-with-substitution.
- New command file `cmd/init.go`.
- `cmd/create.go`: replace direct AGENTS.md/ac.md writes with a `template.Apply(taskTemplateDir, taskDir, vars)` call.
- `cmd/contextfiles.go`: shrinks to hold the *default* content used to prepopulate `task-template/` when missing.
- New directory under the task root — visible to the user, who can edit it freely.
- No effect on existing open tasks unless the user manually re-runs scaffolding (not provided by this change).
