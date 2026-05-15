## Why

`task.json` stores `repo_path` per repository, pointing at the base repo on disk. This duplicates information already in the repo registry and confuses AI agents reading the file: they see a path and assume it's where their work lives, when in fact the worktree is at `repositories/<repo-name>/code/`. The base repo is internal plumbing, not the agent's workspace.

## What Changes

- Remove the `repo_path` field from `task.json` (the `RepoState` struct in `internal/task`).
- `wrt close` and `wrt reopen` now resolve the base repo path by looking up the repo in the registry (by name) instead of reading it from `task.json`.
- `AGENTS.md` gains an explicit instruction telling agents that work happens in `repositories/<repo>/code/` (and `repositories/<repo>/backports/<version>/`), and that the base repo path is registry-internal and should not be touched.
- No backward compatibility shim — task.json files written by older builds will fail to load. Accepted, since this tool has one user.

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `task-lifecycle`: `task.json` schema sheds the redundant `repo_path`; close/reopen rely on the registry as the single source of truth for repo paths.
- `task-context-files`: AGENTS.md gains a workspace-vs-base-repo distinction.

## Impact

- `internal/task/task.go`: remove `RepoPath` from `RepoState`.
- `cmd/repo.go` `runRepoAdd`: stop writing `RepoPath` into the new repo state.
- `cmd/reopen.go`: load the registry entry by repo name; use its `Path` for `git worktree add`.
- `cmd/close.go`: same — registry lookup for `git worktree remove`.
- `cmd/contextfiles.go`: extend the AGENTS.md template with the workspace clarification.
- Existing task.json files written before this change will fail to load. No migration is provided.
