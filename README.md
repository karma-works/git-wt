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
| **List worktrees** | `wtf list` | `git worktree list` |
| **Create worktree** | `wtf create feat` | `git worktree add -b feat ../repo.feat` |
| **Clean up** | `wtf remove` | `cd ../repo && git worktree remove ../repo.feat && git branch -d feat` |
| **List with status** | `wtf list` | `git worktree list` (paths only) |

## Installation

### 1. Install the binary
```bash
go install github.com/hacr/wtf@latest
```
Ensure your `$GOPATH/bin` is in your `PATH`.

## Quick Start (Recommended)

The most efficient way to use `wtf` is via the `exec` command, which handles the entire lifecycle of a task:

```bash
# 1. Create a worktree
# 2. Run the agent
# 3. Prompt for cleanup
wtf exec -b feature-name "Implement the user login flow"
```

## Shell Integration

To enable automatic directory switching for `create` and `remove` commands, add the following to your shell profile.

### Bash / Zsh
Add to `.bashrc` or `.zshrc`:
```bash
# This function makes 'wtf create' and 'wtf remove' 
# automatically change your current directory.
wtf() {
    if [[ "$1" == "list" || "$1" == "exec" ]]; then
        command wtf "$@"
        return
    }
    
    local out
    out=$(command wtf "$@")
    if [[ -d "$out" ]]; then
        cd "$out"
    else
        echo "$out"
    fi
}
```

### PowerShell
Add to `$PROFILE`:
```powershell
function wtf {
    if ($args[0] -eq "list" -or $args[0] -eq "exec") {
        & wtf.exe $args
        return
    }

    $out = & wtf.exe $args
    if (Test-Path $out -PathType Container) {
        Set-Location $out
    } else {
        $out
    }
}
```

## Usage

### List worktrees
```bash
wtf list
```

### Create a new worktree
```bash
# If using the shell integration:
wtf create feature-name

# Otherwise, pipe to cd:
wtf create feature-name | cd
```

### Execute an agent prompt in a new worktree
```bash
# Creates 'feat-logging' worktree, runs 'opencode' agent, then prompts to remove
wtf exec -b feat-logging "Implement the logging logic using logrus"
```
The agent can be configured via `WTF_AGENT` environment variable (defaults to `opencode`).

### Remove worktree
From inside the worktree directory:
```bash
wtf remove | cd
```

Or remove a specific worktree by name with interactive confirmation:
```bash
wtf remove -i feat-logging
```

## License
MIT
