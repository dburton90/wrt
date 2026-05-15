## Why

Navigating to a task directory requires knowing the exact task name and typing `cd $(wrt path <name>)`. A `wrt-cd` shell function with an fzf picker would let users jump to any task with a single command. Task names also currently have no character restrictions, which can break path resolution and the first-word extraction that `wrt-cd` depends on.

## What Changes

- `wrt path` accepts a full `wrt list` output line as its argument (extracts task name from the first word), so it can be used directly with fzf output
- `wrt create` rejects task names containing spaces or slashes
- `wrt completion bash` and `wrt completion zsh` append a `wrt-cd` shell function definition to their output; sourcing the completion script gives users `wrt-cd` for free

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `shell-completion`: completion scripts now also emit the `wrt-cd` shell function
- `task-lifecycle`: `wrt create` gains name validation; `wrt path` gains first-word extraction

## Impact

- `cmd/completion.go`: append `wrt-cd` function after bash/zsh completion output
- `cmd/path.go`: use `strings.Fields(args[0])[0]` to extract task name
- `cmd/create.go`: validate name with `strings.ContainsAny(name, " /")`
- No new dependencies; fzf is an external tool users provide
- Fish completion unchanged (skip for now)
