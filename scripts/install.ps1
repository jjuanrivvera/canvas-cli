# Canvas CLI Installation Script for Windows (PowerShell)

param(
    [string]$InstallDir = "$env:LOCALAPPDATA\Programs\canvas-cli",
    [switch]$AddToPath = $true
)

$ErrorActionPreference = "Stop"

# Configuration
$REPO = "jjuanrivvera/canvas-cli"
$BINARY_NAME = "canvas.exe"

function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Green
}

function Write-Warn {
    param([string]$Message)
    Write-Host "[WARN] $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

function Get-LatestVersion {
    Write-Info "Fetching latest release..."

    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO/releases/latest"
        $version = $response.tag_name
        Write-Info "Latest version: $version"
        return $version
    }
    catch {
        Write-Error "Failed to fetch latest version: $_"
        exit 1
    }
}

function Get-Architecture {
    $arch = [System.Environment]::GetEnvironmentVariable("PROCESSOR_ARCHITECTURE")

    switch ($arch) {
        "AMD64" { return "x86_64" }
        "ARM64" { return "arm64" }
        default {
            Write-Error "Unsupported architecture: $arch"
            exit 1
        }
    }
}

function Install-Binary {
    param(
        [string]$Version,
        [string]$Architecture
    )

    $downloadUrl = "https://github.com/$REPO/releases/download/$Version/canvas_$($Version.TrimStart('v'))_Windows_$Architecture.zip"
    $tempDir = [System.IO.Path]::GetTempPath() + "canvas-cli-" + [System.Guid]::NewGuid().ToString()
    $zipFile = "$tempDir\canvas.zip"

    Write-Info "Creating temporary directory: $tempDir"
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null

    Write-Info "Downloading from: $downloadUrl"

    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $zipFile
    }
    catch {
        Write-Error "Failed to download binary: $_"
        Remove-Item -Path $tempDir -Recurse -Force
        exit 1
    }

    Write-Info "Extracting archive..."
    Expand-Archive -Path $zipFile -DestinationPath $tempDir -Force

    Write-Info "Installing to $InstallDir..."
    if (!(Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

    Copy-Item -Path "$tempDir\$BINARY_NAME" -Destination "$InstallDir\$BINARY_NAME" -Force

    Remove-Item -Path $tempDir -Recurse -Force

    Write-Info "Installed successfully!"
}

function Add-ToPath {
    param([string]$Directory)

    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")

    if ($currentPath -notlike "*$Directory*") {
        Write-Info "Adding $Directory to PATH..."

        $newPath = "$currentPath;$Directory"
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")

        # Update current session
        $env:Path = "$env:Path;$Directory"

        Write-Info "Added to PATH successfully!"
        Write-Warn "You may need to restart your terminal for PATH changes to take effect"
    }
    else {
        Write-Info "$Directory is already in PATH"
    }
}

function Test-Installation {
    Write-Info "Verifying installation..."

    $binaryPath = "$InstallDir\$BINARY_NAME"

    if (!(Test-Path $binaryPath)) {
        Write-Error "Binary not found at $binaryPath"
        exit 1
    }

    try {
        $versionOutput = & $binaryPath version
        Write-Info "Verification successful!"
        Write-Host $versionOutput
    }
    catch {
        Write-Error "Failed to run binary: $_"
        exit 1
    }
}

function Install-Completion {
    Write-Info "Installing PowerShell completion..."

    $profilePath = $PROFILE.CurrentUserAllHosts

    if (!(Test-Path $profilePath)) {
        New-Item -ItemType File -Path $profilePath -Force | Out-Null
    }

    $completionCommand = "Invoke-Expression -Command `"& '$InstallDir\$BINARY_NAME' completion powershell`""

    $profileContent = Get-Content -Path $profilePath -ErrorAction SilentlyContinue

    if ($profileContent -notcontains $completionCommand) {
        Add-Content -Path $profilePath -Value "`n# Canvas CLI completion"
        Add-Content -Path $profilePath -Value $completionCommand

        Write-Info "Added completion to PowerShell profile"
        Write-Warn "Restart PowerShell or run: . `$PROFILE"
    }
    else {
        Write-Info "Completion already configured"
    }
}

# Main installation process
function Main {
    Write-Host "========================================"
    Write-Host "  Canvas CLI Installation Script"
    Write-Host "========================================"
    Write-Host ""

    $architecture = Get-Architecture
    Write-Info "Detected architecture: $architecture"

    $version = Get-LatestVersion

    Install-Binary -Version $version -Architecture $architecture

    if ($AddToPath) {
        Add-ToPath -Directory $InstallDir
    }

    Test-Installation

    Write-Host ""
    $installCompletion = Read-Host "Install PowerShell completion? (Y/n)"

    if ($installCompletion -ne "n" -and $installCompletion -ne "N") {
        Install-Completion
    }

    Write-Host ""
    Write-Host "========================================"
    Write-Info "Installation complete!"
    Write-Host "========================================"
    Write-Host ""
    Write-Host "Next steps:"
    Write-Host "  1. Restart your terminal (for PATH changes)"
    Write-Host "  2. Authenticate: canvas auth login --instance https://canvas.instructure.com"
    Write-Host "  3. Test: canvas courses list"
    Write-Host "  4. Get help: canvas --help"
    Write-Host ""
    Write-Host "Documentation: https://github.com/$REPO/tree/main/docs"
}

Main
