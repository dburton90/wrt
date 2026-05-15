package cmd

import (
	"github.com/dburton90/wrt/internal/config"
	"github.com/dburton90/wrt/internal/task"
	"github.com/spf13/cobra"
)

func loadConfig() (*config.Config, error) {
	return config.Load()
}

func loadTaskRootForCompletion() (string, error) {
	cfg, err := loadConfig()
	if err != nil {
		return "", err
	}
	return cfg.TaskRoot()
}

func completeOpenTaskNames(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	taskRoot, err := loadTaskRootForCompletion()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return task.OpenNames(taskRoot), cobra.ShellCompDirectiveNoFileComp
}

func completeClosedTaskNames(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	taskRoot, err := loadTaskRootForCompletion()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return task.ClosedNames(taskRoot), cobra.ShellCompDirectiveNoFileComp
}

func completeAllTaskNames(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	taskRoot, err := loadTaskRootForCompletion()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	names := append(task.OpenNames(taskRoot), task.ClosedNames(taskRoot)...)
	return names, cobra.ShellCompDirectiveNoFileComp
}
