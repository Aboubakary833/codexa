$version = $args[0] | ForEach-Object { if ($_){$_} else {"latest"} }
$binName = "codexa.exe"
$installDir = "$env:LOCALAPPDATA\Programs\Codexa"
$tmpDir = New-Item -ItemType Directory -Path ([System.IO.Path]::GetTempPath() + [System.Guid]::NewGuid())


New-Item -ItemType Directory -Force -Path $installDir | Out-Null

# Download
if ($version -eq "latest") {
    $url = (Invoke-RestMethod "https://api.github.com/repos/aboubakary833/codexa/releases/latest").assets |
           Where-Object { $_.name -like "*windows_amd64.zip" } |
           Select-Object -ExpandProperty browser_download_url
} else {
    $url = "https://github.com/aboubakary833/codexa/releases/download/$version/codexa_$version_windows_amd64.zip"
}

Write-Host "Downloading Codexa $version..."
Invoke-WebRequest -Uri $url -OutFile "$tmpDir\codexa.zip"

# Extraction
Write-Host "Extracting..."
Expand-Archive -Path "$tmpDir\codexa.zip" -DestinationPath $tmpDir

# Installation
Move-Item -Path "$tmpDir\codexa.exe" -Destination "$installDir\$binName" -Force

# Add codexa to PATH
if (-not ($env:PATH -split ";" | Where-Object { $_ -eq $installDir })) {
    [Environment]::SetEnvironmentVariable("PATH", "$installDir;$env:PATH", "User")
    Write-Host "Added Codexa to PATH. Restart your shell to use it."
}

# Generate codexa PowerShell ompletions
Write-Host "Generating PowerShell completion..."
$completionScript = "$installDir\codexa-completion.ps1"
$null = "$installDir\$binName" completion powershell > $completionScript
if (-not (Select-String -Path $PROFILE -Pattern "codexa-completion.ps1" -Quiet)) {
    Add-Content -Path $PROFILE -Value ". '$completionScript'"
}

Write-Host "Codexa installed successfully! Restart PowerShell."
