package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the task root (creates directories and default template; idempotent)",
	Args:  cobra.NoArgs,
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(_ *cobra.Command, _ []string) error {
	_, taskRoot := mustTaskRoot()
	created, err := initTaskRoot(taskRoot)
	if err != nil {
		return err
	}
	if len(created) == 0 {
		fmt.Printf("Task root at %s is already initialized; nothing to do.\n", taskRoot)
		return nil
	}
	fmt.Printf("Initialized task root at %s\n", taskRoot)
	for _, c := range created {
		fmt.Printf("  created: %s\n", c)
	}
	return nil
}

// initTaskRoot ensures the standard task-root layout exists and prepopulates
// tasks/task-template/ with default content (skip-if-exists per file). Returns
// the list of paths that were created on this call (empty if nothing was new).
//
// Used by `wrt init` and as a lazy fallback by `wrt create`.
func initTaskRoot(taskRoot string) ([]string, error) {
	var created []string

	dirs := []string{
		filepath.Join(taskRoot, "tasks", "open"),
		filepath.Join(taskRoot, "tasks", "closed"),
		filepath.Join(taskRoot, "repos"),
		filepath.Join(taskRoot, "tasks", "task-template"),
	}
	for _, d := range dirs {
		if _, err := os.Stat(d); err == nil {
			continue
		}
		if err := os.MkdirAll(d, 0o755); err != nil {
			return created, fmt.Errorf("creating %s: %w", d, err)
		}
		created = append(created, d)
	}

	templateDir := filepath.Join(taskRoot, "tasks", "task-template")
	for relPath, content := range defaultTemplate {
		dst := filepath.Join(templateDir, relPath)
		if _, err := os.Stat(dst); err == nil {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return created, fmt.Errorf("creating dir for %s: %w", dst, err)
		}
		if err := os.WriteFile(dst, []byte(content), 0o644); err != nil {
			return created, fmt.Errorf("writing %s: %w", dst, err)
		}
		created = append(created, dst)
	}

	return created, nil
}
