# wrt

Task workspace manager with git worktrees. Organizes work on tasks (tickets, issues, or freeform) by managing git worktrees across multiple repositories, with backporting support and AI agent context files.

## Install

**From source:**

```sh
git clone https://github.com/dburton90/wrt
cd wrt
go build -o ~/.local/bin/wrt .
```

**Using `go install`:**

```sh
go install github.com/dburton90/wrt@latest
```

Requires Go 1.21+.

## Setup

Create `~/.config/wrt/config.toml`:

```toml
task_root = "/home/you/tasks"
# username = "yourname"  # optional, defaults to $USER
```

Enable shell completion (zsh example):

```sh
wrt completion zsh > ~/.zsh/completions/_wrt
```

## Quick Start

```sh
# Register a repository
wrt repo create myrepo --path /path/to/repo --base-branch main

# Create a task
wrt create JIRA-111 --url https://jira.example.com/JIRA-111

# Navigate into the task and add the repo (creates worktree + branch)
cd $(wrt path JIRA-111)
wrt repo add myrepo

# Do your work in repositories/myrepo/code/ ...

# Add a backport worktree
wrt backport add myrepo 1.1.1

# Close the task when done (removes worktrees, saves patches)
wrt close JIRA-111

# Reopen later (recreates worktrees)
wrt reopen JIRA-111
```

## Commands

| Command | Description |
|---|---|
| `wrt create [name]` | Create a new task workspace |
| `wrt list` | List open tasks (newest first) |
| `wrt info <task>` | Show task details |
| `wrt path <task>` | Print task directory path |
| `wrt close <task>` | Close task: save patches, remove worktrees |
| `wrt reopen <task>` | Reopen task: recreate worktrees |
| `wrt repo create <name>` | Register a repository |
| `wrt repo list` | List registered repositories |
| `wrt repo add <name>` | Add repo to current task (run inside task folder) |
| `wrt backport add <repo> <version>` | Add backport worktree (run inside task folder) |
| `wrt completion <shell>` | Generate shell completion (bash/zsh/fish) |

## Directory Structure

```
<task-root>/
├── repos/
│   └── <repo-name>/
│       └── repo.json          ← repo registry entry
└── tasks/
    ├── open/
    │   └── <task-id>/
    │       ├── task.json
    │       ├── AGENTS.md      ← auto-generated context for AI agents
    │       ├── ac.md          ← acceptance criteria
    │       └── repositories/
    │           └── <repo>/
    │               ├── code/          ← worktree (task branch)
    │               └── backports/
    │                   └── <version>/ ← worktree (backport branch)
    └── closed/
        └── <task-id>/
            └── repositories/
                └── <repo>/
                    └── <repo>.patch   ← saved patch
```
