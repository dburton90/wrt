package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:       "completion <shell>",
	Short:     "Generate shell completion script",
	ValidArgs: []string{"bash", "zsh", "fish"},
	Args:      cobra.ExactArgs(1),
	RunE:      runCompletion,
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

const wrtCdFunc = `
# wrt-cd: cd to a task directory using fzf
wrt-cd() { local l; l=$(wrt list | fzf --header-lines=2); [ -n "$l" ] && cd "$(wrt path "$l")"; }
`

func runCompletion(_ *cobra.Command, args []string) error {
	switch args[0] {
	case "bash":
		if err := rootCmd.GenBashCompletion(os.Stdout); err != nil {
			return err
		}
		fmt.Fprint(os.Stdout, wrtCdFunc)
		return nil
	case "zsh":
		if err := rootCmd.GenZshCompletion(os.Stdout); err != nil {
			return err
		}
		fmt.Fprint(os.Stdout, wrtCdFunc)
		return nil
	case "fish":
		return rootCmd.GenFishCompletion(os.Stdout, true)
	default:
		return fmt.Errorf("unsupported shell %q. Choose: bash, zsh, fish", args[0])
	}
}
