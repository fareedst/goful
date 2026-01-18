<#
.SYNOPSIS
    Helper script for running Goful in Windows Docker container

.DESCRIPTION
    Builds and runs Goful in a Windows container with proper terminal settings.
    [IMPL:DOCKERFILE_WINDOWS] [ARCH:DOCKER_WINDOWS_BUILD] [REQ:DOCKER_WINDOWS_CONTAINER]

.PARAMETER Build
    Force rebuild of the Docker image even if it already exists.

.PARAMETER Args
    Additional arguments to pass to goful.

.EXAMPLE
    .\docker-run.ps1
    Runs goful interactively in a Windows container.

.EXAMPLE
    .\docker-run.ps1 -Build
    Rebuilds the Docker image and runs goful.

.EXAMPLE
    .\docker-run.ps1 --help
    Passes --help to goful inside the container.
#>

param(
    [switch]$Build,
    [Parameter(ValueFromRemainingArguments=$true)]
    [string[]]$GofulArgs
)

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ImageName = "goful:windows"
$DockerfilePath = Join-Path $ScriptDir "Dockerfile.windows"

# Check if Docker is available
if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
    Write-Error "Docker is not installed or not in PATH"
    exit 1
}

# Check if Docker is in Windows containers mode
$dockerInfo = docker info 2>&1 | Select-String "OSType"
if ($dockerInfo -notmatch "windows") {
    Write-Warning "Docker may not be in Windows containers mode."
    Write-Warning "Switch to Windows containers: Right-click Docker Desktop tray icon > 'Switch to Windows containers...'"
}

# Build image if requested or if it doesn't exist
$ImageExists = $false
try {
    $null = docker image inspect $ImageName 2>&1
    $ImageExists = $true
} catch {
    $ImageExists = $false
}

if ($Build -or -not $ImageExists) {
    Write-Host "Building Docker image: $ImageName" -ForegroundColor Cyan
    docker build -f $DockerfilePath -t $ImageName $ScriptDir
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Docker build failed"
        exit $LASTEXITCODE
    }
}

# Run container interactively
Write-Host "Starting Goful in Windows container..." -ForegroundColor Cyan
docker run -it --rm `
    -v "${ScriptDir}:C:\workspace" `
    -w "C:\workspace" `
    $ImageName @GofulArgs

exit $LASTEXITCODE
