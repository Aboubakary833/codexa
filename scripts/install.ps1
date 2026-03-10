$version = $args[0] | ForEach-Object { if ($_){$_} else {"latest"} }
$binName = "codexa.exe"
$installDir = "$env:LOCALAPPDATA\Programs\Codexa"
$tmpDir = Join-Path ([System.IO.Path]::GetTempPath()) ([System.Guid]::NewGuid().ToString())

New-Item -ItemType Directory -Force -Path $installDir | Out-Null
New-Item -ItemType Directory -Force -Path $tmpDir | Out-Null

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

# Find the extracted binary automatically
$binaryPath = Get-ChildItem -Path $tmpDir -Filter "codexa*.exe" -Recurse | Select-Object -First 1
if (-not $binaryPath) {
    Write-Error "No codexa executable found after extraction"
    exit 1
}

# Installation
Move-Item -Path $binaryPath.FullName -Destination "$installDir\$binName" -Force

# Add codexa to PATH
if (-not ($env:PATH -split ";" | Where-Object { $_ -eq $installDir })) {
    [Environment]::SetEnvironmentVariable("PATH", "$installDir;$env:PATH", "User")
    $env:PATH = "$installDir;$env:PATH"
    Write-Host "Added Codexa to PATH. Restart your shell to use it."
}

# Generate PowerShell completions
Write-Host "Generating PowerShell completion..."
$completionScript = "$installDir\codexa-completion.ps1"
& "$installDir\$binName" completion powershell > $completionScript
if (-not (Select-String -Path $PROFILE -Pattern "codexa-completion.ps1" -Quiet)) {
    Add-Content -Path $PROFILE -Value ". '$completionScript'"
}

Write-Host "Codexa installed successfully! Restart PowerShell."
