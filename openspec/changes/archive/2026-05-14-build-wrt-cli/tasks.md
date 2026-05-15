## 1. Project Setup

- [x] 1.1 Initialize Go module (`go mod init github.com/<user>/wrt`) and create top-level directory structure (`cmd/`, `internal/`, `main.go`)
- [x] 1.2 Add dependencies: Cobra (CLI), promptui (interactive prompts), BurntSushi/toml (config parsing)
- [x] 1.3 Wire up root Cobra command with version flag and stub subcommands

## 2. Config

- [x] 2.1 Implement `internal/config` package: load `~/.config/wrt/config.toml`, apply defaults, expose `TaskRoot()` and `Username()` helpers
- [x] 2.2 Implement username sanitization (strip spaces, lowercase) used for `{user}` substitution
- [x] 2.3 Error with actionable message when `task_root` is not configured

## 3. Repo Registry

- [x] 3.1 Implement `internal/registry` package: `repo.json` struct with JSON marshal/unmarshal, load and save functions
- [x] 3.2 Implement `wrt repo create <name>` — interactive prompts + flags, write `repo.json`, error on duplicate
- [x] 3.3 Implement `wrt repo list` — scan `<task-root>/repos/`, print formatted table

## 4. Task Context Detection

- [x] 4.1 Implement `internal/taskctx` package: walk CWD upward looking for `task.json`, return task directory or error

## 5. Task Lifecycle — Create, List, Info, Path

- [x] 5.1 Implement `internal/task` package: `task.json` struct, load/save, create new task directory with `task.json`
- [x] 5.2 Implement `wrt create` — interactive prompts + flags, validate no name conflict in open/closed, write task directory
- [x] 5.3 Implement `wrt list` — scan `tasks/open/`, sort by `created` desc, print table
- [x] 5.4 Implement `wrt info <task-name>` — load task.json, print details, flag missing worktrees
- [x] 5.5 Implement `wrt path <task-name>` — print task directory path only, exit non-zero if not found

## 6. Task Context Files

- [x] 6.1 Implement AGENTS.md template generation (folder structure, ac.md description, log.md instruction)
- [x] 6.2 Implement ac.md scaffold (header + placeholder content)
- [x] 6.3 Write both files during `wrt create`

## 7. Worktree Management — Add Repo

- [x] 7.1 Implement branch name template substitution (`{user}`, `{task-id}`, `{version}`)
- [x] 7.2 Implement `wrt repo add <name>` — look up registry, create branch from base, run `git worktree add`, update `task.json`
- [x] 7.3 Handle error cases: repo not in registry (hint to `wrt repo list` / `wrt repo create`), already added to task

## 8. Backport Worktrees

- [x] 8.1 Implement `wrt backport add <repo> <version>` — resolve version alias, create backport branch, run `git worktree add`, update `task.json`
- [x] 8.2 Handle error cases: version not in repo.json, repo not in task, already exists

## 9. Task Lifecycle — Close and Reopen

- [x] 9.1 Implement `wrt close <task-name>` — run `git format-patch` per repo, run `git worktree remove`, move task dir to `tasks/closed/`
- [x] 9.2 Implement `wrt reopen <task-name>` — move task dir to `tasks/open/`, recreate worktrees, run `git am`; on conflict leave paused and print instructions

## 10. Shell Completion

- [x] 10.1 Implement `wrt completion <shell>` command (bash/zsh/fish via Cobra)
- [x] 10.2 Add `ValidArgsFunction` for task name arguments (open tasks for most commands, closed tasks for `reopen`)
- [x] 10.3 Add `ValidArgsFunction` for `wrt repo add` (registered repo names)
- [x] 10.4 Add `ValidArgsFunction` for `wrt backport add` first arg (repos in current task) and second arg (backport versions from repo.json)

## 11. Integration and Polish

- [x] 11.1 Add consistent error message format across all commands (stderr, non-zero exit)
- [x] 11.2 Test full happy-path flow end-to-end: init config → repo create → task create → repo add → backport add → close → reopen
- [x] 11.3 Write installation instructions (build from source, `go install`)
