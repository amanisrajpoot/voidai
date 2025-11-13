# PowerShell script for Windows installation

Write-Host "Installing Security Platform Desktop Agent..." -ForegroundColor Green

# Check if running as administrator
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if (-not $isAdmin) {
    Write-Host "This script requires administrator privileges. Please run as administrator." -ForegroundColor Red
    exit 1
}

# Install binary
$installPath = "$env:ProgramFiles\SecurityPlatform\Agent"
New-Item -ItemType Directory -Force -Path $installPath | Out-Null
Copy-Item "security-platform-agent.exe" -Destination "$installPath\security-platform-agent.exe" -Force

# Create config directory
$configPath = "$env:ProgramData\SecurityPlatform"
New-Item -ItemType Directory -Force -Path $configPath | Out-Null
if (Test-Path "config.yaml") {
    Copy-Item "config.yaml" -Destination "$configPath\config.yaml" -Force
}

# Install as Windows Service
$serviceName = "SecurityPlatformAgent"
$serviceDisplayName = "Security Platform Desktop Agent"
$serviceDescription = "Security Platform Desktop Agent for observability and security"

# Check if service already exists
$existingService = Get-Service -Name $serviceName -ErrorAction SilentlyContinue
if ($existingService) {
    Stop-Service -Name $serviceName -Force
    sc.exe delete $serviceName
    Start-Sleep -Seconds 2
}

# Create service
New-Service -Name $serviceName `
    -DisplayName $serviceDisplayName `
    -Description $serviceDescription `
    -BinaryPathName "$installPath\security-platform-agent.exe -service" `
    -StartupType Automatic | Out-Null

# Start service
Start-Service -Name $serviceName

Write-Host "Service installed and started successfully!" -ForegroundColor Green
Write-Host "Service Name: $serviceName" -ForegroundColor Cyan
Write-Host "Config Path: $configPath" -ForegroundColor Cyan
