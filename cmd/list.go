package cmd

import (
	"fmt"
	"sort"

	"github.com/dburton90/wrt/internal/task"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all open tasks",
	Args:  cobra.NoArgs,
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(_ *cobra.Command, _ []string) error {
	_, taskRoot := mustTaskRoot()

	tasks, _, err := task.ListOpen(taskRoot)
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		fmt.Println("No open tasks.")
		return nil
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Created.After(tasks[j].Created)
	})

	fmt.Printf("%-20s  %-12s  %-5s  %s\n", "NAME", "CREATED", "REPOS", "DESCRIPTION")
	fmt.Printf("%-20s  %-12s  %-5s  %s\n", "----", "-------", "-----", "-----------")
	for _, t := range tasks {
		desc := t.Description
		if len(desc) > 50 {
			desc = desc[:47] + "..."
		}
		fmt.Printf("%-20s  %-12s  %-5d  %s\n",
			t.Name,
			t.Created.Format("2006-01-02"),
			len(t.Repositories),
			desc,
		)
	}
	return nil
}
