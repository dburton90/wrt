package gitutil

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

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
