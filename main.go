package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hacr/wtf/internal/worktree"
	"github.com/spf13/cobra"
)

var (
	interactiveFlag bool
	branchFlag      string
	forceFlag       bool
)

var rootCmd = &cobra.Command{
	Use:   "wtf",
	Short: "wtf (Work Tree Flow) is a tool to manage git worktrees efficiently",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List worktree paths",
	RunE: func(cmd *cobra.Command, args []string) error {
		paths, err := worktree.List()
		if err != nil {
			return err
		}
		for _, p := range paths {
			fmt.Println(p)
		}
		return nil
	},
}

// prompt asks the user a y/n question.
func askPrompt(message string, defaultYes bool) bool {
	choices := " [y/N] "
	if defaultYes {
		choices = " [Y/n] "
	}
	fmt.Fprintf(os.Stderr, "%s%s", message, choices)

	var response string
	_, err := fmt.Scanln(&response)
	if err != nil && err.Error() != "unexpected newline" {
		return defaultYes
	}

	response = strings.ToLower(strings.TrimSpace(response))
	if response == "" {
		return defaultYes
	}

	if response == "y" || response == "yes" {
		return true
	}
	return false
}

func ensureWorktree(branch string, create bool) (string, error) {
	repoName, err := worktree.GetRepoName()
	if err != nil {
		return "", err
	}

	path := worktree.GetPath(repoName, branch)

	if create {
		// Check if it already exists before trying to add
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Creating worktree at %s...\n", path)
			if err := worktree.Add(branch, path, true); err != nil {
				return "", err
			}
		}
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("worktree path %s does not exist. Use -c or create to create it", path)
	}

	return absPath, nil
}

func runRemove(currentDir, branch, absMainPath string, interactive bool, force bool) error {
	if interactive {
		if !askPrompt(fmt.Sprintf("Are you sure you want to remove worktree for branch %s?", branch), false) {
			return nil
		}
	}

	fmt.Fprintf(os.Stderr, "Removing worktree for branch %s...\n", branch)

	// Move process to root directory to avoid CWD issues after removal
	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}

	env := []string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + os.Getenv("HOME"),
		"PWD=" + absMainPath,
	}

	// 1. Remove worktree
	worktreeArgs := []string{"worktree", "remove"}
	if force {
		worktreeArgs = append(worktreeArgs, "--force")
	}
	worktreeArgs = append(worktreeArgs, currentDir)

	c := exec.Command("git", worktreeArgs...)
	c.Dir = absMainPath
	c.Env = env
	c.Stdout = os.Stderr
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
	}

	// 2. Delete branch
	branchArgs := []string{"branch"}
	if force {
		branchArgs = append(branchArgs, "-D")
	} else {
		branchArgs = append(branchArgs, "-d")
	}
	branchArgs = append(branchArgs, branch)

	c = exec.Command("git", branchArgs...)
	c.Dir = absMainPath
	c.Env = env
	c.Stdout = os.Stderr
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to delete branch %s: %v\n", branch, err)
	}

	fmt.Println(absMainPath)
	return nil
}

func runCreate(branch string) error {
	path, err := ensureWorktree(branch, true)
	if err != nil {
		return err
	}
	fmt.Println(path)
	return nil
}

var removeCmd = &cobra.Command{
	Use:   "remove [branch]",
	Short: "Remove a worktree (defaults to current)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repoName, err := worktree.GetRepoName()
		if err != nil {
			return err
		}

		var branch string
		var currentDir string

		if len(args) > 0 {
			branch = args[0]
			path := worktree.GetPath(repoName, branch)
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			currentDir = absPath
		} else {
			currentDir, err = os.Getwd()
			if err != nil {
				return err
			}

			branch, err = worktree.GetCurrentBranch()
			if err != nil {
				return err
			}
		}

		if branch == "main" || branch == "master" {
			return fmt.Errorf("cannot remove the main branch/worktree")
		}

		mainRepoPath := filepath.Join("..", repoName)
		absMainPath, err := filepath.Abs(mainRepoPath)
		if err != nil {
			return err
		}

		return runRemove(currentDir, branch, absMainPath, interactiveFlag, forceFlag)
	},
}

var execCmd = &cobra.Command{
	Use:   "exec [prompt]",
	Short: "Execute an agent prompt in a new worktree",
	RunE: func(cmd *cobra.Command, args []string) error {
		if branchFlag == "" {
			return fmt.Errorf("branch (-b) is required")
		}

		path, err := ensureWorktree(branchFlag, true)
		if err != nil {
			return err
		}

		agent := os.Getenv("WTF_AGENT")
		if agent == "" {
			agent = "opencode"
		}

		promptStr := strings.Join(args, " ")
		fmt.Fprintf(os.Stderr, "Executing agent %s with prompt: %s\n", agent, promptStr)

		agentArgs := []string{"--prompt", promptStr}
		c := exec.Command(agent, agentArgs...)
		c.Dir = path
		fmt.Fprintf(os.Stderr, "DEBUG: Executing %s with args %v in dir %s\n", agent, agentArgs, path)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		if err := c.Run(); err != nil {
			return fmt.Errorf("agent failed: %w", err)
		}

		if askPrompt(fmt.Sprintf("Remove worktree %s?", branchFlag), true) {
			repoName, _ := worktree.GetRepoName()
			mainRepoPath := filepath.Join("..", repoName)
			absMainPath, _ := filepath.Abs(mainRepoPath)
			return runRemove(path, branchFlag, absMainPath, false, false)
		}

		fmt.Println(path)
		return nil
	},
}

func init() {
	removeCmd.Flags().BoolVarP(&interactiveFlag, "interactive", "i", false, "Prompt for confirmation before removal")
	removeCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force removal of worktree and branch")
	execCmd.Flags().StringVarP(&branchFlag, "branch", "b", "", "Branch name to create/use")

	// Add 'create' command
	createCmd := &cobra.Command{
		Use:   "create [branch]",
		Short: "Create and switch to a new worktree",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(args[0])
		},
	}

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(execCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
