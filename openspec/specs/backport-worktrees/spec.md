## ADDED Requirements

### Requirement: Add backport worktree
`wrt backport add <repo-name> <version>` SHALL look up the version alias in the repo's `backport_branches` map to resolve the base branch, create a new git branch using the backport branch template, and create a git worktree at `repositories/<repo-name>/backports/<version>/`. It SHALL update `task.json` to record the version in the repo's `backport_branches` list.

#### Scenario: Valid backport add
- **WHEN** `wrt backport add some-repo 1.1.1` is run inside task `JIRA-111`
- **THEN** branch `users/dbarton/JIRA-111-backport-1.1.1` is created from `release/1.1.1`, worktree is created at `repositories/some-repo/backports/1.1.1/`, and `task.json` is updated

#### Scenario: Version not in repo's backport_branches
- **WHEN** `wrt backport add some-repo 5.0.x` is run and `5.0.x` is not in `repo.json`
- **THEN** the tool errors: "version '5.0.x' is not configured for repo 'some-repo'. Run `wrt repo list` to see available backport versions."

#### Scenario: Repo not added to task
- **WHEN** `wrt backport add some-repo 1.1.1` is run but some-repo is not in `task.json`
- **THEN** the tool errors: "repo 'some-repo' is not part of this task. Run `wrt repo add some-repo` first."

#### Scenario: Backport version already added
- **WHEN** `wrt backport add some-repo 1.1.1` is run and that backport already exists in `task.json`
- **THEN** the tool errors: "backport '1.1.1' for repo 'some-repo' already exists in this task"

### Requirement: Backport worktree path
Backport worktrees SHALL be created at `<task-dir>/repositories/<repo-name>/backports/<version>/`.

#### Scenario: Backport path
- **WHEN** backport `1.1.1` is added for repo `some-repo` in task `JIRA-111`
- **THEN** the worktree path is `tasks/open/JIRA-111/repositories/some-repo/backports/1.1.1/`

### Requirement: Backport branch naming
The backport branch name SHALL be resolved from `backport_branch_template` in `repo.json` with `{user}`, `{task-id}`, and `{version}` substituted. The default template is `users/{user}/{task-id}-backport-{version}`.

#### Scenario: Default template
- **WHEN** `repo.json` uses the default backport branch template
- **THEN** `wrt backport add some-repo 1.1.1` for task `JIRA-111` by user `dbarton` creates branch `users/dbarton/JIRA-111-backport-1.1.1`

#### Scenario: Custom template
- **WHEN** `repo.json` sets `backport_branch_template` to `backport/{task-id}/{version}`
- **THEN** the created branch is `backport/JIRA-111/1.1.1`
