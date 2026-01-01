function wt {
    param(
        [Parameter(ValueFromRemainingArguments = $true)]
        $Arguments
    )

    # Create a temporary file for directives
    $directiveFile = [System.IO.Path]::GetTempFileName()
    
    # Set the environment variable for the gitwt process
    $env:GITWT_DIRECTIVE_FILE = $directiveFile
    
    try {
        # Run gitwt
        & gitwt @Arguments
        
        # If the directive file has content, execute it
        if (Test-Path $directiveFile) {
            $content = Get-Content $directiveFile
            if ($content) {
                # In PowerShell, we can simply execute the commands
                # We assume git-wt writes valid PowerShell commands (like cd "path")
                # Though currently it writes POSIX style 'cd "path"', which PS handles fine
                Invoke-Expression $content
            }
        }
    }
    finally {
        # Cleanup
        if (Test-Path $directiveFile) {
            Remove-Item $directiveFile
        }
        Remove-Item Env:\GITWT_DIRECTIVE_FILE
    }
}
