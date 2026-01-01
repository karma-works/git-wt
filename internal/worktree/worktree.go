package worktree

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// List returns a list of worktree paths.
func List() ([]string, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %w", err)
	}

	var paths []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "worktree ") {
			paths = append(paths, strings.TrimPrefix(line, "worktree "))
		}
	}
	return paths, nil
}

// Add creates a new worktree.
func Add(branch, path string, create bool) error {
	args := []string{"worktree", "add"}
	if create {
		args = append(args, "-b", branch)
		args = append(args, path)
	} else {
		args = append(args, path, branch)
	}

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Remove deletes a worktree and its associated branch.
func Remove(path string, branch string) error {
	// 1. Remove worktree
	cmd := exec.Command("git", "worktree", "remove", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
	}

	// 2. Delete branch
	if branch != "" {
		cmd = exec.Command("git", "branch", "-d", branch)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: failed to delete branch %s: %v\n", branch, err)
		}
	}

	return nil
}

// GetRepoName returns the name of the current repository.
func GetRepoName() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get repo root: %w", err)
	}
	return filepath.Base(strings.TrimSpace(string(output))), nil
}

// GetPath returns the conventional worktree path for a branch.
func GetPath(repoName, branch string) string {
	return filepath.Join("..", fmt.Sprintf("%s.%s", repoName, branch))
}

// GetCurrentBranch returns the current branch name.
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}
