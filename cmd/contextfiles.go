package cmd

// defaultTemplate maps paths within tasks/task-template/ to their default content.
// Used to prepopulate the template directory on first init. Files already present
// are not overwritten.
var defaultTemplate = map[string]string{
	"AGENTS.md": `# Task Workspace: {task-id}

This directory is a **wrt** task workspace. It contains git worktrees and context files for working on this task across one or more repositories.

## Directory Structure

` + "```" + `
repositories/
└── <repo-name>/
    ├── code/           ← main task worktree
    └── backports/
        └── <version>/ ← backport worktree
` + "```" + `

Each worktree is a full git checkout. You can run git commands, edit files, and run tests directly inside them.

## Where to Work

**All work happens inside the per-repo worktrees** under ` + "`repositories/<repo>/code/`" + ` (and ` + "`repositories/<repo>/backports/<version>/`" + ` for backports). These are the only places you should read, edit, and commit code for this task.

The "base repository" that backs each worktree lives outside this task directory and is registered in wrt's internal repo registry. It is **wrt's plumbing — do not modify it directly** (do not ` + "`cd`" + ` into it, do not commit there, do not run git operations against it). Treat it as read-only infrastructure.

## Context Files

- **ac.md** — Acceptance criteria for this task. Describes what must be true for the task to be considered complete. Read this before starting work.
- **log.md** — Agent activity log. Record all significant actions here (see instructions below).

## Agent Instructions

When working in this workspace, record all significant actions in **log.md** using this format:

` + "```" + `
[YYYY-MM-DD HH:MM] <what was done and why>
` + "```" + `

Examples:
- ` + "`" + `[2026-05-14 10:30] Created implementation plan based on ac.md requirements` + "`" + `
- ` + "`" + `[2026-05-14 11:15] Modified auth.go — ac.md requires tokens expire after 24h` + "`" + `
- ` + "`" + `[2026-05-14 14:00] Ran unit tests: 3 failing, investigating token refresh path` + "`" + `

Log entries help track progress and make it easy to resume work after interruptions.
`,

	"ac.md": `# Acceptance Criteria

<!-- List the criteria that must be met for this task to be considered complete. -->
<!-- Each criterion should be verifiable — you know when it's done. -->

- [ ]
`,

	".claude/settings.local.json": `{
  "permissions": {
    "allow": [
      "Bash(cd repositories/*)",
      "Bash(cd {task-dir}/repositories/*)"
    ]
  }
}
`,
}
