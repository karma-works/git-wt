# wtf - Work Tree Flow

`wtf` (Work Tree Flow) is a Go-based wrapper for `git worktree` designed to streamline your workflow, especially when working with coding agents like Claude or GitHub Copilot.

## Features

- **Standardized Paths**: Automatically manages worktrees in a consistent `../<repo>.<branch>` structure.
- **Easy Switching**: Create and enter worktrees with a single command.
- **Agent Integration**: Start your favorite AI agent immediately upon creating a worktree with the `-x` flag.
- **Clean Removal**: Removes the worktree and its associated branch in one step.
- **Shell Support**: Includes wrappers for Bash, Zsh, and PowerShell to enable automatic directory switching.

## Comparison: wtf vs Plain Git

| Task | wtf | Plain Git |
| :--- | :--- | :--- |
| **Switch worktrees** | `wtf switch feat` | `cd ../repo.feat` |
| **Create + Start Claude** | `wtf switch -c -x claude feat` | `git worktree add -b feat ../repo.feat && cd ../repo.feat && claude` |
| **Create + Start Antigravity** | `wtf switch -c -x "antigravity chat 'Fix bug'" feat` | `git worktree add -b feat ../repo.feat && cd ../repo.feat && antigravity chat 'Fix bug'` |
| **Create + Open VS Code** | `wtf switch -c -x "code ." feat` | `git worktree add -b feat ../repo.feat && cd ../repo.feat && code .` |
| **Clean up** | `wtf remove` | `cd ../repo && git worktree remove ../repo.feat && git branch -d feat` |
| **List with status** | `wtf list` | `git worktree list` (paths only) |

## Installation

### 1. Install the binary
```bash
go install github.com/hacr/wtf@latest
```
Ensure your `$GOPATH/bin` is in your `PATH`.

### 2. Set up Shell Wrapper (Recommended)

Since a binary cannot change the current directory of the parent shell, use the provided wrappers.

#### Bash / Zsh
Add the following to your `.bashrc` or `.zshrc`:
```bash
# Source the wrapper script or define as a function
wtf() {
    # Replace with actual path to the wtf wrapper
    source /path/to/wtf "$@"
}
```

#### PowerShell
Add the following to your `$PROFILE`:
```powershell
# Source the wrapper script
. /path/to/wtf/wtf.ps1
```

## Usage

### List worktrees
```bash
wtf list
```

### Create and switch to a new branch/worktree
```bash
wtf switch -c feature-name
```

### Create and start an agent
```bash
# Start Claude
wtf switch -c -x "claude" feat-agent

# Start an Antigravity session with a specific prompt
wtf switch -c -x "antigravity chat 'Implement the logger'" feat-logging

# Open in VS Code immediately
wtf switch -c -x "code ." feat-ui

# Start GitHub Copilot CLI session
wtf switch -c -x "gh copilot suggest 'Write a python script'" feat-copilot
```

### Remove current worktree
From inside the worktree directory:
```bash
wtf remove
```

## License
MIT
