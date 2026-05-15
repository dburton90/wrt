## Context

`wrt` is a greenfield Go CLI with no existing codebase to migrate. The tool operates entirely on the local filesystem and delegates all git operations to the `git` binary on PATH. It has no server component, no database, and no network requirements beyond what `git` itself uses.

## Goals / Non-Goals

**Goals:**
- Single binary, no runtime dependencies beyond `git`
- All state stored as human-readable files (JSON, Markdown)
- Works from any directory via CWD traversal
- Shell autocomplete for all arguments

**Non-Goals:**
- Syncing state across machines
- Integrating with Jira/GitHub APIs (task IDs are just strings)
- Running or orchestrating agents (wrt prepares context; agents run separately)
- Managing git operations beyond worktree create/delete and format-patch

## Decisions

### Go + Cobra for CLI framework

**Decision**: Use Go with [Cobra](https://github.com/spf13/cobra) and [promptui](https://github.com/manifoldco/promptui) for interactive prompts.

**Rationale**: Cobra provides shell completion generation (bash/zsh/fish) for free via `RegisterFlagCompletionFunc` and `ValidArgsFunction`. Single binary distribution. Strong standard library for filesystem and exec operations.

**Alternatives considered**: Python (Click) — good ergonomics but requires runtime install. Node — too heavy for a system tool.

### JSON for task.json and repo.json

**Decision**: Store `task.json` and `repo.json` as JSON, not TOML or YAML.

**Rationale**: JSON is unambiguous, has no indentation sensitivity, and is trivially parseable by agents and scripts. Human-editability is secondary — the CLI is the intended write path.

**Alternatives considered**: TOML — good for config but less standard for structured data with maps. YAML — too many footguns.

### TOML for `~/.config/wrt/config.toml`

**Decision**: Global config uses TOML (via `github.com/BurntSushi/toml`).

**Rationale**: TOML is idiomatic for Go CLI config files. The config file is human-written (not machine-generated), so readability matters here.

### CWD traversal for task context

**Decision**: Context-sensitive commands (`wrt repo add`, `wrt backport add`) find the task root by walking up from CWD looking for `task.json`, identical to how `git` finds `.git`.

**Rationale**: Allows running commands from anywhere inside the task workspace (e.g., deep inside `repositories/some-repo/code/src/`). Fail with a clear error if no `task.json` is found in any ancestor.

### git format-patch for close/reopen

**Decision**: `wrt close` runs `git format-patch <base-branch>..HEAD -o <patch-dir>/` for each repo worktree. `wrt reopen` runs `git am <patch-file>` on a fresh worktree.

**Rationale**: Portable, human-readable, and survives the deletion of branches. On conflict, `git am` leaves the worktree in a paused state — wrt detects this, informs the user, and exits without aborting, so the user can resolve manually.

**Alternatives considered**: `git bundle` — binary format, harder to inspect. Copying the worktree directory — huge disk usage.

### Branch naming via templates in repo.json

**Decision**: Branch templates are stored in `repo.json` with `{user}`, `{task-id}`, and `{version}` as substitution tokens. Defaults: `users/{user}/{task-id}` and `users/{user}/{task-id}-backport-{version}`. `{user}` resolves to `$USER` at runtime, overridable in `config.toml`.

**Rationale**: Different repos in the same org may enforce different branch naming conventions. Per-repo templates handle this without requiring a global override mechanism.

### Task root directory layout

```
<task-root>/
├── repos/
│   └── <repo-name>/
│       └── repo.json
└── tasks/
    ├── open/
    │   └── <task-id>/
    │       ├── task.json
    │       ├── AGENTS.md
    │       ├── ac.md
    │       └── repositories/
    │           └── <repo-name>/
    │               ├── code/          ← git worktree
    │               └── backports/
    │                   └── <version>/ ← git worktree
    └── closed/
        └── <task-id>/
            ├── task.json
            ├── AGENTS.md
            ├── ac.md
            └── repositories/
                └── <repo-name>/
                    └── <repo-name>.patch
```

**Rationale**: Separating `repos/` (global registry) from `tasks/` (per-task workspaces) makes the registry independently manageable. `open/` and `closed/` separation allows `wrt list` to scan only open tasks efficiently.

## Risks / Trade-offs

- **Patch conflicts on reopen** → Mitigation: Leave `git am` in paused state, print clear instructions for `git am --continue` or `git am --abort`. Never auto-abort.
- **Worktree path length** → Deep nesting (`tasks/open/JIRA-111/repositories/some-repo/code/`) can hit filesystem path limits on some systems. Mitigation: document the limitation; no automatic fix.
- **Stale task.json after manual git ops** → If a user manually deletes a branch that's recorded in task.json, `wrt info` may show stale data. Mitigation: `wrt info` verifies worktree existence and flags discrepancies.
- **$USER with spaces** → Unlikely but possible. Mitigation: strip spaces and lowercase when substituting `{user}` in branch templates.
