package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const DefaultTaskBranchTemplate = "users/{user}/{task-id}"
const DefaultBackportBranchTemplate = "users/{user}/{task-id}-backport-{version}"

type Repo struct {
	Name                   string            `json:"name"`
	Path                   string            `json:"path"`
	Link                   string            `json:"link,omitempty"`
	DefaultBaseBranch      string            `json:"default_base_branch"`
	TaskBranchTemplate     string            `json:"task_branch_template"`
	BackportBranchTemplate string            `json:"backport_branch_template"`
	BackportBranches       map[string]string `json:"backport_branches,omitempty"`
}

func RepoDir(taskRoot, name string) string {
	return filepath.Join(taskRoot, "repos", name)
}

func RepoPath(taskRoot, name string) string {
	return filepath.Join(RepoDir(taskRoot, name), "repo.json")
}

func Load(taskRoot, name string) (*Repo, error) {
	p := RepoPath(taskRoot, name)
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf(
				"repo %q not found\n\nRun `wrt repo list` to see registered repos or `wrt repo create %s` to add it",
				name, name,
			)
		}
		return nil, fmt.Errorf("reading repo: %w", err)
	}
	var r Repo
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("parsing repo.json: %w", err)
	}
	if r.Path == "" {
		return nil, fmt.Errorf("repo.json for %q is missing required field: path", name)
	}
	if r.DefaultBaseBranch == "" {
		return nil, fmt.Errorf("repo.json for %q is missing required field: default_base_branch", name)
	}
	return &r, nil
}

func Save(taskRoot string, r *Repo) error {
	if r.TaskBranchTemplate == "" {
		r.TaskBranchTemplate = DefaultTaskBranchTemplate
	}
	if r.BackportBranchTemplate == "" {
		r.BackportBranchTemplate = DefaultBackportBranchTemplate
	}
	dir := RepoDir(taskRoot, r.Name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating repo dir: %w", err)
	}
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(RepoPath(taskRoot, r.Name), data, 0o644)
}

func List(taskRoot string) ([]*Repo, error) {
	reposDir := filepath.Join(taskRoot, "repos")
	entries, err := os.ReadDir(reposDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading repos dir: %w", err)
	}
	var repos []*Repo
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		r, err := Load(taskRoot, e.Name())
		if err != nil {
			continue
		}
		repos = append(repos, r)
	}
	return repos, nil
}

func Names(taskRoot string) []string {
	repos, _ := List(taskRoot)
	names := make([]string, 0, len(repos))
	for _, r := range repos {
		names = append(names, r.Name)
	}
	return names
}
