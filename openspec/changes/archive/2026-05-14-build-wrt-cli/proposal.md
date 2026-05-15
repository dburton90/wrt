## Why

Developers working on multi-repo projects with backporting requirements spend significant time manually managing git worktrees, tracking task context across repositories, and preparing structured context for AI agents. `wrt` eliminates this overhead by providing a single CLI that owns the full lifecycle of a task workspace.

## What Changes

- New Go CLI tool `wrt` installable as a single binary
- Global task root directory (user-configured) stores all task workspaces and repo registry
- Task workspaces contain git worktrees, structured context files, and agent instructions
- Repo registry maps repo names to paths, branch templates, and backport branch mappings
- Task lifecycle: create → add repos → add backports → close (saves patches) → reopen (restores from patches)
- Auto-generated `AGENTS.md` gives AI agents structured context about the workspace

## Capabilities

### New Capabilities

- `config`: Global configuration management — task root path, username
- `repo-registry`: Register and manage known repositories with their branch templates and backport mappings
- `task-lifecycle`: Create, list, inspect, close, and reopen tasks
- `worktree-management`: Add repository worktrees to tasks, create and track task branches
- `backport-worktrees`: Add backport worktrees to tasks after implementation, branch naming with version suffixes
- `task-context-files`: Auto-generate `AGENTS.md`, scaffold `ac.md` at task creation
- `shell-completion`: Tab completion for task names, repo names, and backport versions

### Modified Capabilities

## Impact

- New Go module and binary, no existing code affected
- Requires `git` on PATH for worktree operations
- Writes to user's filesystem under the configured task root
- Shell completion requires one-time setup per shell (bash/zsh/fish)
