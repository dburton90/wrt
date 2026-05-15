package taskctx

import (
	"fmt"
	"os"
	"path/filepath"
)

// Find walks up from CWD looking for task.json, returning the directory that contains it.
func Find() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting working directory: %w", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "task.json")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf(
		"not inside a task directory\n\nNavigate into a task folder or run `wrt list` to find one",
	)
}
