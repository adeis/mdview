Write-Host "==============================================" -ForegroundColor Blue
Write-Host " Installing mdview CLI globally on Windows... " -ForegroundColor Blue
Write-Host "==============================================" -ForegroundColor Blue

# Check Wails
$wails = Get-Command wails -ErrorAction SilentlyContinue
if (-not $wails) {
    $gopath = go env GOPATH
    if (Test-Path "$gopath\bin\wails.exe") {
        $wails = "$gopath\bin\wails.exe"
    } else {
        Write-Host "Error: Wails CLI is not installed." -ForegroundColor Red
        Write-Host "Please run: go install github.com/wailsapp/wails/v2/cmd/wails@latest" -ForegroundColor Yellow
        exit 1
    }
}

# 1. Build application
Write-Host "`n[1/3] Building the application..." -ForegroundColor Yellow
& $wails build
if ($LASTEXITCODE -ne 0) {
    Write-Host "Error: Build failed." -ForegroundColor Red
    exit 1
}

# 2. Check binary
$srcBin = "build\bin\mdviewer.exe"
if (-not (Test-Path $srcBin)) {
    Write-Host "Error: Compiled binary not found." -ForegroundColor Red
    exit 1
}

# 3. Copy to WindowsApps folder (which is in user PATH by default)
Write-Host "`n[2/3] Installing binary to user PATH..." -ForegroundColor Yellow
$targetDir = "$env:LOCALAPPDATA\Microsoft\WindowsApps"
$targetPath = "$targetDir\mdview.exe"

Copy-Item -Path $srcBin -Destination $targetPath -Force

if (Test-Path $targetPath) {
    Write-Host "`n[3/3] Installation successful!" -ForegroundColor Green
    Write-Host "You can now run: mdview <file.md> from any Command Prompt / PowerShell!" -ForegroundColor Blue
} else {
    Write-Host "Error: Failed to install to $targetPath. Check folder permissions." -ForegroundColor Red
}
Write-Host "==============================================" -ForegroundColor Blue
