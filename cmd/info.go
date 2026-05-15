package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dburton90/wrt/internal/task"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:               "info <task-name>",
	Short:             "Show task details",
	Args:              cobra.ExactArgs(1),
	RunE:              runInfo,
	ValidArgsFunction: completeOpenTaskNames,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfo(_ *cobra.Command, args []string) error {
	_, taskRoot := mustTaskRoot()
	name := args[0]

	taskDir, open, err := task.Find(taskRoot, name)
	if err != nil {
		return err
	}

	t, err := task.Load(taskDir)
	if err != nil {
		return err
	}

	status := "open"
	if !open {
		status = "closed"
	}

	fmt.Printf("Task:    %s [%s]\n", t.Name, status)
	if t.URL != "" {
		fmt.Printf("URL:     %s\n", t.URL)
	}
	if t.Description != "" {
		fmt.Printf("Desc:    %s\n", t.Description)
	}
	fmt.Printf("Created: %s\n", t.Created.Format("2006-01-02 15:04 UTC"))
	fmt.Printf("Path:    %s\n", taskDir)

	if len(t.Repositories) == 0 {
		fmt.Println("\nNo repositories added. Run `wrt repo add <name>` inside the task folder.")
		return nil
	}

	fmt.Println("\nRepositories:")
	for repoName, rs := range t.Repositories {
		worktreePath := filepath.Join(taskDir, "repositories", repoName, "code")
		worktreeStatus := "ok"
		if open {
			if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
				worktreeStatus = "⚠ worktree missing"
			}
		}
		fmt.Printf("  %s\n", repoName)
		fmt.Printf("    branch:  %s [%s]\n", rs.TaskBranch, worktreeStatus)
		fmt.Printf("    base:    %s\n", rs.BaseBranch)

		for _, bp := range rs.BackportBranches {
			bpPath := filepath.Join(taskDir, "repositories", repoName, "backports", bp.Version)
			bpStatus := "ok"
			if open {
				if _, err := os.Stat(bpPath); os.IsNotExist(err) {
					bpStatus = "⚠ worktree missing"
				}
			}
			fmt.Printf("    backport %s: %s [%s]\n", bp.Version, bp.Branch, bpStatus)
		}
	}
	return nil
}
