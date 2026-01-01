#!/bin/bash

# wt - a wrapper for gitwt to handle directory switching using a directive file
# Usage: source wt.sh [args]
# Or add this function to your .bashrc / .zshrc:
# wt() { source /path/to/wt.sh "$@"; }

wt() {
    # Create a temporary file for directives
    local directive_file=$(mktemp)
    
    # Export it so gitwt can see it
    export GITWT_DIRECTIVE_FILE="$directive_file"
    
    # Run gitwt (stdour/stderr go to terminal as usual)
    gitwt "$@"
    
    # If the directive file has content, source it
    if [ -s "$directive_file" ]; then
        source "$directive_file"
    fi
    
    # Cleanup
    rm -f "$directive_file"
    unset GITWT_DIRECTIVE_FILE
}

# If the script is being sourced, define the function
if [[ "${BASH_SOURCE[0]}" != "${0}" ]] || [[ -n "$ZSH_VERSION" ]]; then
    # Function already defined above
    :
else
    # If run as a script, we warn the user
    echo "Warning: wt.sh should be sourced or used as a function to enable directory switching."
    wt "$@"
fi