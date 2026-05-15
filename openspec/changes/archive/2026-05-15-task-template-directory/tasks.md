## 1. Template engine

- [x] 1.1 Create `internal/template` package
- [x] 1.2 Implement `Apply(srcDir, dstDir string, vars map[string]string) error` that walks `srcDir`, mirrors structure to `dstDir`, substitutes vars in file contents via `strings.NewReplacer`, skips destination files that already exist
- [x] 1.3 Variables: `{task-id}`, `{task-dir}`, `{task-root}`, `{user}`
- [x] 1.4 Substitution is content-only; file and directory names are not substituted

## 2. Default template payload

- [x] 2.1 In `cmd/contextfiles.go` (or new file): define the default template content as in-memory `map[string]string` (relative path → content)
- [x] 2.2 AGENTS.md default: replace `fmt.Sprintf` with embedded literal using `{task-id}`; include the worktree-vs-base-repo clarification (depends on change `strip-repo-path` or duplicates that wording)
- [x] 2.3 ac.md default: identical to today's `writeAcMD` content
- [x] 2.4 `.claude/settings.local.json` default: minimal allowlist with relative + `{task-dir}` absolute entries

## 3. Bootstrap helper

- [x] 3.1 Implement `initTaskRoot(taskRoot string) error`: ensures `tasks/open`, `tasks/closed`, `repos`, `tasks/task-template` exist; writes any missing default template files (skip-if-exists)
- [x] 3.2 Used by both `wrt init` and the lazy fallback in `wrt create`

## 4. `wrt init` command

- [x] 4.1 Create `cmd/init.go` with a `wrt init` Cobra command
- [x] 4.2 Loads config, resolves task root, errors if `task_root` is unset
- [x] 4.3 Calls `initTaskRoot`; prints what was created and what already existed
- [x] 4.4 Idempotent: a second run reports "already initialized"

## 5. `wrt create` integration

- [x] 5.1 In `cmd/create.go`, after task directory is created, call `initTaskRoot` to ensure templates exist
- [x] 5.2 Replace direct AGENTS.md/ac.md writes with `template.Apply(<task-root>/tasks/task-template, taskDir, vars)`
- [x] 5.3 Remove or repurpose `writeAgentsMD`/`writeAcMD` (their content now lives only as defaults inside the bootstrap helper)
- [x] 5.4 Verify the same files exist post-create as before, plus `.claude/settings.local.json`

## 6. Manual verification

- [x] 6.1 `wrt init` on a fresh task root creates the four subdirs and four default template files
- [x] 6.2 `wrt init` re-run is idempotent and reports no changes
- [x] 6.3 `wrt create JIRA-999` on a task root WITHOUT a template directory: directory is created on the fly, scaffolding works
- [x] 6.4 `wrt create JIRA-999` on a task root WITH custom additions in `task-template/`: those files land in the task with substitution applied
- [x] 6.5 `{task-dir}` in the .claude allowlist contains the correct absolute task path
- [x] 6.6 Hand-edit a template file, create another task: edit is reflected
- [x] 6.7 Skip-if-exists: pre-create `ac.md` manually in destination → run `wrt create` → existing content preserved
