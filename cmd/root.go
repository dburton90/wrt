package cmd

import (
	"fmt"
	"os"

	"github.com/dburton90/wrt/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wrt",
	Short: "Task workspace manager with git worktrees",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = "0.1.0"
	rootCmd.SilenceUsage = true
}

// mustConfig loads config and exits on error.
func mustConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		fatal("%s", err)
	}
	return cfg
}

// mustTaskRoot loads config and returns task root, exiting if not configured.
func mustTaskRoot() (cfg *config.Config, taskRoot string) {
	cfg = mustConfig()
	root, err := cfg.TaskRoot()
	if err != nil {
		fatal("%s", err)
	}
	return cfg, root
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
