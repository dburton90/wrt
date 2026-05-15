## Context

The repo registry at `<task-root>/repos/<repo>/repo.json` already authoritatively records every repo's on-disk path. Today, `wrt repo add` copies that path into `task.json` as `repo_path` so that `close`/`reopen` can run git commands without a registry lookup. The cost of that duplication: AI agents reading `task.json` see a base-repo path and conclude it's their working directory, then try to commit there — pointing at the *original* repo rather than the worktree.

The fix is small: drop the duplicate. The registry is always reachable (it lives next to `task.json` in the task root), so the cost of resolving the path on demand is a single file read.

## Goals / Non-Goals

**Goals:**
- `task.json` no longer contains `repo_path`.
- `wrt close` and `wrt reopen` continue to function by resolving the base repo path through the registry.
- AGENTS.md tells agents unambiguously where to work.

**Non-Goals:**
- Migration of existing task.json files. Old tasks will break.
- Restructuring the registry. It stays as-is.
- Generalizing AGENTS.md beyond the workspace clarification (template-based scaffolding lives in a separate change).

## Decisions

### Registry lookup in close and reopen

`close.go` and `reopen.go` currently iterate `task.Repositories` and use `rs.RepoPath` directly. After this change, they call `registry.Load(taskRoot, repoName)` to retrieve the path. The registry name is the map key already, so no schema change is needed there. An additional file read per repo on close/reopen is negligible.

**Alternative considered**: cache the path on the in-memory `Task` struct after a registry pass. Rejected — solves a non-problem (close/reopen are interactive, one-time operations).

### AGENTS.md wording

The structure section already lists `repositories/<repo>/code/` as the layout. We add an explicit "Agent instructions" sentence: code lives in the worktree at `repositories/<repo>/code/`; the base repo (in the registry) is wrt's plumbing and SHALL NOT be modified directly.

## Risks / Trade-offs

- **Old task.json files break.** Anyone with open tasks at the time of upgrade must close or recreate them. Acceptable for a single-user tool.
- **Registry drift.** If a repo is renamed or its entry is deleted between `repo add` and `close`, the lookup fails. Previously, the cached `repo_path` would have survived a registry deletion. Net positive though — the registry IS the source of truth, so divergence should surface.
- **More registry reads.** Trivial — close/reopen are rare interactive commands.
