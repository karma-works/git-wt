function wtf {
    param(
        [Parameter(ValueFromRemainingArguments = $true)]
        $Arguments
    )

    # Create a temporary file for directives
    $directiveFile = [System.IO.Path]::GetTempFileName()
    
    # Set the environment variable for the wtf process
    $env:WTF_DIRECTIVE_FILE = $directiveFile
    
    try {
        # Run wtf binary
        & wtf.exe @Arguments
        
        # If the directive file has content, execute it
        if (Test-Path $directiveFile) {
            $content = Get-Content $directiveFile
            if ($content) {
                Invoke-Expression $content
            }
        }
    }
    finally {
        # Cleanup
        if (Test-Path $directiveFile) {
            Remove-Item $directiveFile
        }
        Remove-Item Env:\WTF_DIRECTIVE_FILE
    }
}
