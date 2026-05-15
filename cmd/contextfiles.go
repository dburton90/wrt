package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func writeAgentsMD(taskDir, taskName string) error {
	content := fmt.Sprintf(`# Task Workspace: %s

This directory is a **wrt** task workspace. It contains git worktrees and context files for working on this task across one or more repositories.

## Directory Structure

`+"```"+`
repositories/
└── <repo-name>/
    ├── code/           ← main task worktree
    └── backports/
        └── <version>/ ← backport worktree
`+"```"+`

Each worktree is a full git checkout. You can run git commands, edit files, and run tests directly inside them.

## Context Files

- **ac.md** — Acceptance criteria for this task. Describes what must be true for the task to be considered complete. Read this before starting work.
- **log.md** — Agent activity log. Record all significant actions here (see instructions below).

## Agent Instructions

When working in this workspace, record all significant actions in **log.md** using this format:

`+"```"+`
[YYYY-MM-DD HH:MM] <what was done and why>
`+"```"+`

Examples:
- `+"`"+`[2026-05-14 10:30] Created implementation plan based on ac.md requirements`+"`"+`
- `+"`"+`[2026-05-14 11:15] Modified auth.go — ac.md requires tokens expire after 24h`+"`"+`
- `+"`"+`[2026-05-14 14:00] Ran unit tests: 3 failing, investigating token refresh path`+"`"+`

Log entries help track progress and make it easy to resume work after interruptions.
`, taskName)
	return os.WriteFile(filepath.Join(taskDir, "AGENTS.md"), []byte(content), 0o644)
}

func writeAcMD(taskDir string) error {
	content := `# Acceptance Criteria

<!-- List the criteria that must be met for this task to be considered complete. -->
<!-- Each criterion should be verifiable — you know when it's done. -->

- [ ]
`
	return os.WriteFile(filepath.Join(taskDir, "ac.md"), []byte(content), 0o644)
}
