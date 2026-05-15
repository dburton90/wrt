## Context

`wrt path <task-name>` already prints a task directory path suitable for use in `cd $(wrt path ...)`. The missing piece is a quick way to select a task without knowing its name upfront. `wrt list` already outputs tasks sorted newest-first with a 2-line header; fzf handles interactive selection from that list. The shell function `wrt-cd` ties these together and lives in the completion script so users get it for free when they source completion.

Task names are currently unrestricted strings used directly as directory names, which can produce broken paths (spaces) or unexpected nesting (slashes). Validation is cheap to add at creation time.

## Goals / Non-Goals

**Goals:**
- `wrt-cd` shell function available after sourcing `wrt completion bash|zsh`
- `wrt path` accepts a full `wrt list` line or a bare task name interchangeably
- `wrt create` rejects names with spaces or slashes

**Non-Goals:**
- Fish shell support (deferred)
- Fuzzy matching built into `wrt` itself (fzf handles this externally)
- Validating existing task names retroactively

## Decisions

### First-word extraction in `wrt path`

`wrt path` receives whatever fzf outputs after the user selects a line — a full formatted line like `task-foo   2025-01-10   2   fix the login bug`. Rather than adding a flag or a separate command, `runPath` uses `strings.Fields(args[0])[0]` to extract the first word. Since task names are validated to contain no spaces, the first word is always the name. Bare names (no extra fields) continue to work unchanged.

**Alternative considered**: separate `wrt cd` subcommand. Rejected — unnecessary indirection when `wrt path` already does the job.

### Shell function in completion output

`completion.go` appends the `wrt-cd` function after the cobra-generated script for bash and zsh. The function body is identical for both shells:

```bash
wrt-cd() { local l; l=$(wrt list | fzf --header-lines=2); [ -n "$l" ] && cd "$(wrt path "$l")"; }
```

`--header-lines=2` keeps the `NAME / ----` header visible but non-selectable. The `[ -n "$l" ]` guard means ctrl-c or empty selection is a no-op.

**Alternative considered**: a standalone script users must source separately. Rejected — sourcing completion is already the standard setup step.

### Name validation placement

Validation runs in `runCreate` immediately after the name is determined (whether from CLI arg or interactive prompt), before the conflict check. A single `strings.ContainsAny(name, " /")` call covers both cases with one code path.

## Risks / Trade-offs

- **fzf not installed** → `wrt list | fzf` fails with a clear shell error; `wrt path` is unaffected. No silent failure.
- **Long task names** → fzf truncates display naturally; path output is unaffected.
- **`wrt path ""` (empty first word)** → `strings.Fields("")` returns an empty slice; `[0]` would panic. Guard: if `args[0]` is blank after field-splitting, return an error before indexing.
