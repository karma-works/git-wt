# git-wt

`git-wt` is a Go-based wrapper for `git worktree` designed to streamline your workflow, especially when working with coding agents like Claude or GitHub Copilot.

## Features

- **Standardized Paths**: Automatically manages worktrees in a consistent `../<repo>.<branch>` structure.
- **Easy Switching**: Create and enter worktrees with a single command.
- **Agent Integration**: Start your favorite AI agent immediately upon creating a worktree with the `-x` flag.
- **Clean Removal**: Removes the worktree and its associated branch in one step.
- **Shell Support**: Includes wrappers for Bash, Zsh, and PowerShell to enable automatic directory switching.

## Comparison: git-wt vs Plain Git

| Task | git-wt | Plain Git |
| :--- | :--- | :--- |
| **Switch worktrees** | `wt switch feat` | `cd ../repo.feat` |
| **Create + Start Claude** | `wt switch -c -x claude feat` | `git worktree add -b feat ../repo.feat && cd ../repo.feat && claude` |
| **Create + Start Antigravity** | `wt switch -c -x "antigravity chat 'Fix bug'" feat` | `git worktree add -b feat ../repo.feat && cd ../repo.feat && antigravity chat 'Fix bug'` |
| **Create + Open VS Code** | `wt switch -c -x "code ." feat` | `git worktree add -b feat ../repo.feat && cd ../repo.feat && code .` |
| **Clean up** | `wt remove` | `cd ../repo && git worktree remove ../repo.feat && git branch -d feat` |
| **List with status** | `wt list` | `git worktree list` (paths only) |

## Installation

### 1. Install the binary
```bash
go install github.com/hacr/git-wt@latest
```
Ensure your `$GOPATH/bin` is in your `PATH`.

### 2. Set up Shell Wrapper (Recommended)

Since a binary cannot change the current directory of the parent shell, use the provided wrappers.

#### Bash / Zsh
Add the following to your `.bashrc` or `.zshrc`:
```bash
# Source the wrapper script or define as a function
wt() {
    # Replace with actual path to wt.sh
    source /path/to/git-wt/wt.sh "$@"
}
```

#### PowerShell
Add the following to your `$PROFILE`:
```powershell
# Source the wrapper script
. /path/to/git-wt/git-wt.ps1
```

## Usage

### List worktrees
```bash
wt list
```

### Create and switch to a new branch/worktree
```bash
wt switch -c feature-name
```

### Create and start an agent
```bash
# Start Claude
wt switch -c -x "claude" feat-agent

# Start an Antigravity session with a specific prompt
wt switch -c -x "antigravity chat 'Implement the logger'" feat-logging

# Open in VS Code immediately
wt switch -c -x "code ." feat-ui

# Start GitHub Copilot CLI session
wt switch -c -x "gh copilot suggest 'Write a python script'" feat-copilot
```

### Remove current worktree
From inside the worktree directory:
```bash
wt remove
```

## License
MIT
