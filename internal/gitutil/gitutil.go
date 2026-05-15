package gitutil

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ErrNoUpstream is returned by Upstream when the branch has no upstream tracking ref.
var ErrNoUpstream = errors.New("branch has no upstream tracking ref")

// run executes a git command rooted at dir and returns combined output on error.
func run(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git %v: %w\n%s", args, err, bytes.TrimSpace(out))
	}
	return nil
}

// WorktreeAdd creates a new branch and worktree at destPath rooted from repoPath.
func WorktreeAdd(repoPath, destPath, branch, baseBranch string) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return fmt.Errorf("creating worktree parent dir: %w", err)
	}
	return run(repoPath, "worktree", "add", "-b", branch, destPath, baseBranch)
}

// WorktreeAddExisting creates a worktree for a branch that already exists in the repo.
func WorktreeAddExisting(repoPath, destPath, branch string) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return fmt.Errorf("creating worktree parent dir: %w", err)
	}
	return run(repoPath, "worktree", "add", destPath, branch)
}

// WorktreeRemove removes a worktree from the repo.
func WorktreeRemove(repoPath, worktreePath string) error {
	return run(repoPath, "worktree", "remove", "--force", worktreePath)
}

// FormatPatch writes patches from baseBranch..HEAD in worktreePath to patchFile.
// If there are no commits ahead of baseBranch, the patch file is not created.
func FormatPatch(worktreePath, baseBranch, patchFile string) error {
	cmd := exec.Command("git", "format-patch", baseBranch+"..HEAD", "--stdout")
	cmd.Dir = worktreePath
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("git format-patch: %w", err)
	}
	if len(bytes.TrimSpace(out)) == 0 {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(patchFile), 0o755); err != nil {
		return err
	}
	return os.WriteFile(patchFile, out, 0o644)
}

// ApplyPatch applies a patch file using git am in worktreePath.
// Returns a PatchConflictError if git am exits with conflict.
func ApplyPatch(worktreePath, patchFile string) error {
	data, err := os.ReadFile(patchFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return nil
	}
	cmd := exec.Command("git", "am")
	cmd.Dir = worktreePath
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return &PatchConflictError{Output: string(out)}
	}
	return nil
}

type PatchConflictError struct {
	Output string
}

func (e *PatchConflictError) Error() string {
	return e.Output
}

// IsClean reports whether the worktree has no uncommitted changes (tracked or untracked).
// On false, returns the list of dirty file paths parsed from `git status --porcelain`.
func IsClean(worktreePath string) (bool, []string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = worktreePath
	out, err := cmd.Output()
	if err != nil {
		return false, nil, fmt.Errorf("git status in %s: %w", worktreePath, err)
	}
	trimmed := bytes.TrimSpace(out)
	if len(trimmed) == 0 {
		return true, nil, nil
	}
	var files []string
	for _, line := range bytes.Split(trimmed, []byte("\n")) {
		if len(line) > 3 {
			files = append(files, string(line[3:]))
		}
	}
	return false, files, nil
}

// Upstream resolves the upstream tracking ref of branch in repoPath using
// `git rev-parse --abbrev-ref <branch>@{upstream}`. Returns the remote and the
// remote-side ref name (e.g. remote="origin", ref="main"). When the branch has
// no upstream configured, returns ErrNoUpstream.
func Upstream(repoPath, branch string) (remote, ref string, err error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", branch+"@{upstream}")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return "", "", ErrNoUpstream
	}
	full := strings.TrimSpace(string(out))
	idx := strings.IndexByte(full, '/')
	if idx <= 0 {
		return "", "", fmt.Errorf("unexpected upstream format %q for %s", full, branch)
	}
	return full[:idx], full[idx+1:], nil
}

// Fetch runs `git fetch <remote>` in the base repo.
func Fetch(repoPath, remote string) error {
	return run(repoPath, "fetch", remote)
}

// Rebase runs `git rebase <ontoRef>` in the worktree. On non-zero exit, the
// underlying error includes the git output; callers can probe ConflictedFiles
// to distinguish conflict from other failures.
func Rebase(worktreePath, ontoRef string) error {
	return run(worktreePath, "rebase", ontoRef)
}

// ConflictedFiles returns paths of files with unresolved merge conflicts in worktreePath.
func ConflictedFiles(worktreePath string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", "--diff-filter=U")
	cmd.Dir = worktreePath
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git diff in %s: %w", worktreePath, err)
	}
	trimmed := bytes.TrimSpace(out)
	if len(trimmed) == 0 {
		return nil, nil
	}
	var files []string
	for _, line := range bytes.Split(trimmed, []byte("\n")) {
		files = append(files, string(line))
	}
	return files, nil
}

// HasCommitsAhead returns true if the branch in worktreePath has commits ahead of baseBranch.
func HasCommitsAhead(worktreePath, baseBranch string) bool {
	cmd := exec.Command("git", "log", "--oneline", baseBranch+"..HEAD")
	cmd.Dir = worktreePath
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(bytes.TrimSpace(out)) > 0
}
