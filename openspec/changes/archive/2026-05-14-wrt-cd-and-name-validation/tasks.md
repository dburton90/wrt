## 1. Task Name Validation

- [x] 1.1 In `cmd/create.go`, after name is determined (arg or prompt), add `strings.ContainsAny(name, " /")` check and return an error if true
- [x] 1.2 Verify error message is clear and consistent with existing error style

## 2. First-Word Extraction in wrt path

- [x] 2.1 In `cmd/path.go` `runPath`, replace `name := args[0]` with `strings.Fields(args[0])[0]`, guarding against an empty slice (return error if no fields)
- [x] 2.2 Update `pathCmd` arg validation from `cobra.ExactArgs(1)` to `cobra.ExactArgs(1)` (unchanged — still requires exactly one arg, just parsed differently)

## 3. wrt-cd Shell Function in Completion Output

- [x] 3.1 In `cmd/completion.go`, after `rootCmd.GenBashCompletion(os.Stdout)`, print the `wrt-cd` function definition
- [x] 3.2 In `cmd/completion.go`, after `rootCmd.GenZshCompletion(os.Stdout)`, print the same `wrt-cd` function definition
- [x] 3.3 Verify fish case is left unchanged

## 4. Manual Verification

- [x] 4.1 `wrt create "bad name"` errors with a clear message
- [x] 4.2 `wrt create "bad/name"` errors with a clear message
- [x] 4.3 `wrt path "task-foo   2025-01-10   2   some desc"` prints correct path
- [x] 4.4 `wrt path task-foo` still works as before
- [x] 4.5 `wrt completion zsh` output includes `wrt-cd` function at the end
- [x] 4.6 Sourcing completion and running `wrt-cd` opens fzf picker and cds on selection
