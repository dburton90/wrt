package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dburton90/wrt/internal/task"
	"github.com/dburton90/wrt/internal/template"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [task-name]",
	Short: "Create a new task workspace",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runCreate,
}

var createFlags struct {
	name        string
	url         string
	description string
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&createFlags.url, "url", "", "URL to issue tracker (optional)")
	createCmd.Flags().StringVar(&createFlags.description, "desc", "", "Short description (optional)")
}

func runCreate(cmd *cobra.Command, args []string) error {
	cfg, taskRoot := mustTaskRoot()

	var name string
	if len(args) > 0 {
		name = args[0]
	}

	if name == "" {
		var err error
		name, err = promptString("Task name", "", true)
		if err != nil {
			return err
		}
	}

	if strings.ContainsAny(name, " /") {
		return fmt.Errorf("task name must not contain spaces or slashes")
	}

	url := createFlags.url
	if url == "" && !cmd.Flags().Changed("url") {
		var err error
		url, err = promptOptional("URL")
		if err != nil {
			return err
		}
	}

	description := createFlags.description
	if description == "" && !cmd.Flags().Changed("desc") {
		var err error
		description, err = promptOptional("Description")
		if err != nil {
			return err
		}
	}

	// Check for name conflict in open and closed
	if _, err := os.Stat(filepath.Join(taskRoot, "tasks", "open", name, "task.json")); err == nil {
		return fmt.Errorf("task %q already exists (open). Use `wrt info %s` to inspect it.", name, name)
	}
	if _, err := os.Stat(filepath.Join(taskRoot, "tasks", "closed", name, "task.json")); err == nil {
		return fmt.Errorf("task %q already exists (closed). Use `wrt reopen %s` to reopen it.", name, name)
	}

	// Ensure task root layout and default template exist.
	if _, err := initTaskRoot(taskRoot); err != nil {
		return fmt.Errorf("initializing task root: %w", err)
	}

	taskDir := task.OpenDir(taskRoot, name)
	t := &task.Task{
		Name:        name,
		URL:         url,
		Description: description,
	}
	if err := task.Create(taskDir, t); err != nil {
		return fmt.Errorf("creating task: %w", err)
	}

	vars := map[string]string{
		"task-id":   name,
		"task-dir":  taskDir,
		"task-root": taskRoot,
		"user":      cfg.Username(),
	}
	templateDir := filepath.Join(taskRoot, "tasks", "task-template")
	if err := template.Apply(templateDir, taskDir, vars); err != nil {
		return fmt.Errorf("applying task template: %w", err)
	}

	fmt.Printf("Created task %q at %s\n", name, taskDir)
	return nil
}
