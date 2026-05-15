## Purpose
Maintains the registry of known source repositories at `<task-root>/repos/`, including their on-disk paths, base branches, and branch-name templates.

## Requirements

### Requirement: Repo registry location
The tool SHALL store the repo registry at `<task-root>/repos/<repo-name>/repo.json`. Each registered repo occupies its own subdirectory.

#### Scenario: Registry lookup
- **WHEN** a command references a repo by name (e.g. `some-repo`)
- **THEN** the tool reads `<task-root>/repos/some-repo/repo.json`

### Requirement: repo.json schema
Each `repo.json` SHALL contain:
- `name` (string, required): repo identifier used in CLI commands
- `path` (string, required): absolute path to the git repository on disk
- `link` (string, optional): URL to the remote repository (e.g. GitHub)
- `default_base_branch` (string, required): branch new task branches are created from
- `task_branch_template` (string, required): template for task branch names; default `users/{user}/{task-id}`
- `backport_branch_template` (string, required): template for backport branch names; default `users/{user}/{task-id}-backport-{version}`
- `backport_branches` (object, optional): map of version alias to actual branch name (e.g. `{"1.1.1": "release/1.1.1"}`)

#### Scenario: Valid repo.json
- **WHEN** `repo.json` contains all required fields
- **THEN** the tool can use the repo for worktree operations

#### Scenario: Missing required field
- **WHEN** `repo.json` is missing `path` or `default_base_branch`
- **THEN** the tool errors with a message identifying the missing field

### Requirement: Create repo entry
`wrt repo create <name>` SHALL create a new entry in the repo registry. The command SHALL support both interactive prompts and CLI flags for all fields. It SHALL error if a repo with that name already exists.

#### Scenario: Interactive creation
- **WHEN** `wrt repo create some-repo` is run with no flags
- **THEN** the tool prompts for: path (required), base branch (required), link (optional), branch templates (with defaults shown), backport branches (optional, repeatable)

#### Scenario: Flag-based creation
- **WHEN** `wrt repo create some-repo --path /repos/some-repo --base-branch main` is run
- **THEN** the tool creates `repo.json` without prompting

#### Scenario: Duplicate repo name
- **WHEN** `wrt repo create some-repo` is run and `repos/some-repo/repo.json` already exists
- **THEN** the tool errors: "repo 'some-repo' already exists. Use `wrt repo edit some-repo` to modify it."

### Requirement: List repos
`wrt repo list` SHALL list all registered repos showing: name, path, base branch, and available backport versions.

#### Scenario: No repos registered
- **WHEN** the repos directory is empty
- **THEN** the tool prints "No repos registered. Run `wrt repo create <name>` to add one."

#### Scenario: Repos exist
- **WHEN** one or more repos are registered
- **THEN** the tool prints a formatted list with name, path, base branch, and backport versions for each
