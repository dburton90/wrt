## Context

Today, the only per-task scaffolded files (AGENTS.md, ac.md) are generated from string literals in Go. Adding a third file (the user wants `.claude/settings.local.json` to pre-authorize common agent permissions for the task) reinforces an anti-pattern: every new scaffolded artifact would require modifying the wrt source.

The proposed mechanism is the simplest thing that scales: a directory under the task root whose files are copied into each new task, with templating. It mirrors how `cookiecutter`, `create-next-app`, and many other scaffolders work, scoped to wrt's needs.

The `.claude/settings.local.json` use case forces substitution: Claude Code matches Bash permissions by literal string, so an allowlist entry like `Bash(cd <absolute-task-path>/repositories/*)` must contain the actual path for the specific task. That path is only known at create time. Hence `{task-dir}`.

## Goals / Non-Goals

**Goals:**
- Files under `<task-root>/tasks/task-template/` are copied into every new task created via `wrt create`.
- Substitution of `{task-id}`, `{task-dir}`, `{task-root}`, `{user}` in file *contents* (not in paths).
- `wrt init` is an explicit, idempotent setup command.
- `wrt create` falls back to populating `task-template/` if missing, so the tool works out of the box.
- Default template prepopulation includes AGENTS.md, ac.md, and `.claude/settings.local.json` with sensible defaults.

**Non-Goals:**
- Retroactive scaffolding of existing open tasks.
- Substitution in filenames or directory names.
- Per-task overrides (e.g. `--template-dir` flag). Out of scope for v1.
- Multiple named templates (e.g. `tasks/templates/{minimal,full}/`). One template directory, period.
- Conditional content (Go templates, partials). Plain text-substitution only.

## Decisions

### Substitution syntax: single-brace variables

Use `{task-id}`, `{task-dir}`, `{task-root}`, `{user}` — same form as the existing `task.ResolveBranch` (`{user}`, `{task-id}`, `{version}`). Consistency over invention. Implementation is a `strings.NewReplacer` pass per file, identical in spirit to `ResolveBranch`.

**Alternative considered**: Go's `text/template` package. Rejected — overkill for paths and IDs; introduces a syntax users would have to learn, and brittle interaction with files like `.claude/settings.local.json` that may contain `{...}` JSON braces. (Mitigation possible, but added complexity for no value at this stage.)

### Variable surface

| Variable | Substituted with |
|---|---|
| `{task-id}` | The task name (e.g. `JIRA-111`) |
| `{task-dir}` | Absolute path to the task directory |
| `{task-root}` | Absolute path to the task root |
| `{user}` | Resolved username (sanitized, same as branch-template `{user}`) |

`{created}` (timestamp) is tempting but adds complexity for unclear value. Defer until requested.

### Location: `<task-root>/tasks/task-template/`

User explicitly preferred this over a system-wide `~/.config/wrt/templates/`. Reasons that hold up:
- Per-task-root scoping. A user with separate work/personal task roots gets independent templates.
- Discoverability. `ls <task-root>/tasks/` shows it next to `open/` and `closed/`.
- Self-contained task root. Backing up or moving the task root preserves the templates.

The directory is *not* a task (no `task.json` inside). `wrt list` and friends ignore directories without `task.json`, so this is fine without code changes.

**Alternative considered**: `<task-root>/task-template/` (top-level, sibling of `tasks/` and `repos/`). Slightly cleaner semantic separation, but the user picked the `tasks/` placement and the reasoning above supports it.

### Bootstrap: both `wrt init` AND lazy in `wrt create`

`wrt init`:
- Reads task root from config; errors if not set.
- Creates `tasks/open/`, `tasks/closed/`, `repos/`, `tasks/task-template/` if missing.
- Writes default content into `tasks/task-template/` if any file is missing (skip-if-exists per file).
- Idempotent: running twice is a no-op.

`wrt create` (lazy fallback):
- Before applying the template, checks `tasks/task-template/` exists. If not, calls the same initialization helper that `wrt init` uses to create and prepopulate it. Then proceeds with copy + substitution.

This means a user who never runs `wrt init` still gets templates the first time they `wrt create`. A user who runs `wrt init` first gets to inspect/edit the templates before creating tasks. Both are fine.

### Skip-if-exists, never overwrite

When `wrt create` copies a template file into the task directory, it skips destination paths that already exist. This guarantees the operation is non-destructive in any conceivable race (e.g. user creates task partially by hand and then runs `wrt create`).

### Default content

The defaults baked into the binary, written when prepopulating `task-template/`:

- **`AGENTS.md`** — same content as the current `writeAgentsMD` output, with `{task-id}` substitution restored (today it's `fmt.Sprintf("# Task Workspace: %s", taskName)`; in the template that becomes `# Task Workspace: {task-id}`). Includes the worktree-vs-base-repo clarification from change `strip-repo-path`.
- **`ac.md`** — same as today's `writeAcMD` output. No substitution needed.
- **`.claude/settings.local.json`** — minimal allowlist using `{task-dir}` for absolute path matches:
  ```json
  {
    "permissions": {
      "allow": [
        "Bash(cd repositories/*)",
        "Bash(cd {task-dir}/repositories/*)"
      ]
    }
  }
  ```
  The user can extend this with project-specific entries (jira MCP, opsx skills, openspec) by editing `tasks/task-template/.claude/settings.local.json` once.

### Removal of hardcoded scaffolding

After this change, `cmd/contextfiles.go` no longer exposes `writeAgentsMD` and `writeAcMD` as the path used by `wrt create`. Their content survives as the default template payload, but the writing-into-a-task code path is gone — replaced by the generic copy-with-substitution mechanism.

## Risks / Trade-offs

- **Template breakage.** A user who hand-edits `task-template/` and introduces unbalanced braces or invalid JSON gets that broken file copied verbatim into every new task. Acceptable — same as them writing the broken file directly. Not wrt's job to validate.
- **`{...}` collisions in JSON/code files.** A literal `{user}` in a settings file would be substituted. The default `.claude/settings.local.json` deliberately uses `{task-dir}` for the path. If a future template needs literal `{var}` content, we'd need an escape mechanism. Defer until it actually bites.
- **Existing tasks miss out.** Anyone with open tasks at the time of this change won't get `.claude/settings.local.json`. Workaround: copy it in manually, or close and recreate the task. Acceptable for a single-user tool.
- **`wrt create` behavior changes.** Today, only AGENTS.md and ac.md are written. After this change, anything in `task-template/` (which the user controls) lands in the task. If a user adds a file that conflicts with something `wrt` later adds (e.g. `log.md`), the skip-if-exists rule favors the template — which is consistent with "user owns content".
- **`{task-dir}` in non-portable strings.** Absolute paths in scaffolded files become invalid if the task root moves. A side effect of substituting at create time. Mitigation if needed: relative-path entries (`Bash(cd repositories/*)`) which the default template already includes alongside the absolute form.
