## 1. Schema change

- [x] 1.1 Remove `RepoPath` field from `RepoState` in `internal/task/task.go`
- [x] 1.2 Remove the JSON tag `json:"repo_path"`

## 2. Adjust callers

- [x] 2.1 In `cmd/repo.go` `runRepoAdd`, drop `RepoPath: r.Path` from the `task.RepoState{}` literal
- [x] 2.2 In `cmd/close.go`, replace `rs.RepoPath` with a registry lookup: `r, err := registry.Load(taskRoot, repoName); ...; r.Path`
- [x] 2.3 In `cmd/reopen.go`, do the same: lookup registry by repo name, use `r.Path` for `WorktreeAddExisting`
- [x] 2.4 Surface a clear error if the registry lookup fails ("repo '<name>' no longer registered — cannot operate on this task")

## 3. AGENTS.md content

- [x] 3.1 In `cmd/contextfiles.go`, extend the AGENTS.md template with an explicit instruction: work happens in `repositories/<repo>/code/` (and `…/backports/<version>/`); the base repository on disk is wrt's plumbing and SHALL NOT be modified directly
- [x] 3.2 Verify the wording is unambiguous when read in isolation by an agent

## 4. Manual verification

- [x] 4.1 `wrt repo add` writes a task.json with no `repo_path` key
- [x] 4.2 `wrt close <task>` succeeds without `repo_path` in task.json
- [x] 4.3 `wrt reopen <task>` succeeds without `repo_path` in task.json
- [x] 4.4 `wrt close <task>` errors clearly if the repo is missing from the registry
- [x] 4.5 Generated AGENTS.md contains the workspace-vs-base-repo instruction
