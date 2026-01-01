package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hacr/git-wt/internal/worktree"
	"github.com/spf13/cobra"
)

var (
	createFlag bool
	execFlag   string
)

var rootCmd = &cobra.Command{
	Use:   "gitwt",
	Short: "git-wt is a wrapper to ease the use of git worktree",
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

// writeDirective writes a shell command to the file specified by GITWT_DIRECTIVE_FILE
func writeDirective(command string) {
	directiveFile := os.Getenv("GITWT_DIRECTIVE_FILE")
	if directiveFile == "" {
		return
	}

	f, err := os.OpenFile(directiveFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to write to directive file: %v\n", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(command + "\n"); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to write to directive file: %v\n", err)
	}
}

var switchCmd = &cobra.Command{
	Use:   "switch [branch]",
	Short: "Switch to or create a worktree",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		branch := args[0]
		repoName, err := worktree.GetRepoName()
		if err != nil {
			return err
		}

		path := worktree.GetPath(repoName, branch)

		if createFlag {
			fmt.Printf("Creating worktree at %s...\n", path)
			if err := worktree.Add(branch, path, true); err != nil {
				return err
			}
		}

		// Check if path exists
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			return fmt.Errorf("worktree path %s does not exist. Use -c to create it", path)
		}

		if execFlag != "" {
			fmt.Printf("Executing %s in %s...\n", execFlag, path)
			c := exec.Command("sh", "-c", execFlag)
			c.Dir = absPath
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin
			return c.Run()
		}

		writeDirective(fmt.Sprintf("cd %q", absPath))
		fmt.Printf("Worktree ready. Run: cd %s\n", path)
		return nil
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove the current worktree",
	RunE: func(cmd *cobra.Command, args []string) error {
		repoName, err := worktree.GetRepoName()
		if err != nil {
			return err
		}

		currentDir, err := os.Getwd()
		if err != nil {
			return err
		}

		branch, err := worktree.GetCurrentBranch()
		if err != nil {
			return err
		}

		if branch == "main" || branch == "master" {
			return fmt.Errorf("cannot remove the main branch/worktree")
		}

		// Calculate main repo path
		mainRepoPath := filepath.Join("..", repoName)
		absMainPath, err := filepath.Abs(mainRepoPath)
		if err != nil {
			return err
		}

		fmt.Printf("Removing worktree for branch %s...\n", branch)

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
		c := exec.Command("git", "worktree", "remove", currentDir)
		c.Dir = absMainPath
		c.Env = env
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return fmt.Errorf("failed to remove worktree: %w", err)
		}

		// 2. Delete branch
		c = exec.Command("git", "branch", "-d", branch)
		c.Dir = absMainPath
		c.Env = env
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			fmt.Printf("Warning: failed to delete branch %s: %v\n", branch, err)
		}

		writeDirective(fmt.Sprintf("cd %q", absMainPath))
		fmt.Printf("Worktree removed. Please run: cd %s\n", absMainPath)
		return nil
	},
}

func init() {
	switchCmd.Flags().BoolVarP(&createFlag, "create", "c", false, "Create a new worktree")
	switchCmd.Flags().StringVarP(&execFlag, "exec", "x", "", "Execute a command in the worktree")

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(switchCmd)
	rootCmd.AddCommand(removeCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
