package cmd

import (
	"fmt"
	"strings"

	"github.com/dburton90/wrt/internal/task"
	"github.com/spf13/cobra"
)

var pathCmd = &cobra.Command{
	Use:               "path <task-name>",
	Short:             "Print the path to a task directory",
	Args:              cobra.ExactArgs(1),
	RunE:              runPath,
	ValidArgsFunction: completeAllTaskNames,
}

func init() {
	rootCmd.AddCommand(pathCmd)
}

func runPath(_ *cobra.Command, args []string) error {
	_, taskRoot := mustTaskRoot()
	fields := strings.Fields(args[0])
	if len(fields) == 0 {
		return fmt.Errorf("task name is empty")
	}
	name := fields[0]

	taskDir, _, err := task.Find(taskRoot, name)
	if err != nil {
		return err
	}
	fmt.Print(taskDir)
	return nil
}
