package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type BackportState struct {
	Version    string `json:"version"`
	BaseBranch string `json:"base_branch"`
	Branch     string `json:"branch"`
}

type RepoState struct {
	BaseBranch       string          `json:"base_branch"`
	TaskBranch       string          `json:"task_branch"`
	BackportBranches []BackportState `json:"backport_branches"`
}

type Task struct {
	Name         string               `json:"name"`
	URL          string               `json:"url,omitempty"`
	Description  string               `json:"description,omitempty"`
	Created      time.Time            `json:"created"`
	Repositories map[string]RepoState `json:"repositories"`
}

func OpenDir(taskRoot, name string) string {
	return filepath.Join(taskRoot, "tasks", "open", name)
}

func ClosedDir(taskRoot, name string) string {
	return filepath.Join(taskRoot, "tasks", "closed", name)
}

func JSONPath(taskDir string) string {
	return filepath.Join(taskDir, "task.json")
}

func Load(taskDir string) (*Task, error) {
	data, err := os.ReadFile(JSONPath(taskDir))
	if err != nil {
		return nil, fmt.Errorf("reading task.json: %w", err)
	}
	var t Task
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("parsing task.json: %w", err)
	}
	return &t, nil
}

func Save(taskDir string, t *Task) error {
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(JSONPath(taskDir), data, 0o644)
}

func Create(taskDir string, t *Task) error {
	if t.Created.IsZero() {
		t.Created = time.Now().UTC()
	}
	if t.Repositories == nil {
		t.Repositories = map[string]RepoState{}
	}
	if err := os.MkdirAll(filepath.Join(taskDir, "repositories"), 0o755); err != nil {
		return fmt.Errorf("creating task dir: %w", err)
	}
	return Save(taskDir, t)
}

// Find looks for a task by name in open and closed directories.
// Returns the task directory path and whether the task is open.
func Find(taskRoot, name string) (dir string, open bool, err error) {
	openDir := OpenDir(taskRoot, name)
	if _, err := os.Stat(JSONPath(openDir)); err == nil {
		return openDir, true, nil
	}
	closedDir := ClosedDir(taskRoot, name)
	if _, err := os.Stat(JSONPath(closedDir)); err == nil {
		return closedDir, false, nil
	}
	return "", false, fmt.Errorf("task %q not found", name)
}

func ListOpen(taskRoot string) ([]*Task, []string, error) {
	openDir := filepath.Join(taskRoot, "tasks", "open")
	entries, err := os.ReadDir(openDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		}
		return nil, nil, fmt.Errorf("reading tasks dir: %w", err)
	}
	var tasks []*Task
	var dirs []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(openDir, e.Name())
		t, err := Load(dir)
		if err != nil {
			continue
		}
		tasks = append(tasks, t)
		dirs = append(dirs, dir)
	}
	return tasks, dirs, nil
}

func OpenNames(taskRoot string) []string {
	tasks, _, _ := ListOpen(taskRoot)
	names := make([]string, 0, len(tasks))
	for _, t := range tasks {
		names = append(names, t.Name)
	}
	return names
}

func ClosedNames(taskRoot string) []string {
	closedDir := filepath.Join(taskRoot, "tasks", "closed")
	entries, _ := os.ReadDir(closedDir)
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names
}

// ResolveBranch substitutes {user}, {task-id}, {version} in a template.
func ResolveBranch(template, user, taskID, version string) string {
	r := strings.NewReplacer(
		"{user}", user,
		"{task-id}", taskID,
		"{version}", version,
	)
	return r.Replace(template)
}
