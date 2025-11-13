# PowerShell script to install Security Platform Agent as Windows Service

$serviceName = "SecurityPlatformAgent"
$displayName = "Security Platform Agent"
$description = "Security observability agent for Windows"
$binaryPath = "$PSScriptRoot\security-platform-agent.exe"
$configPath = "$env:ProgramData\SecurityPlatform\config.yaml"

# Check if running as administrator
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $isAdmin) {
    Write-Error "This script must be run as Administrator"
    exit 1
}

# Create service
$service = Get-Service -Name $serviceName -ErrorAction SilentlyContinue
if ($service) {
    Write-Host "Service already exists. Stopping and removing..."
    Stop-Service -Name $serviceName -Force -ErrorAction SilentlyContinue
    sc.exe delete $serviceName
    Start-Sleep -Seconds 2
}

Write-Host "Creating Windows Service..."
New-Service -Name $serviceName `
    -DisplayName $displayName `
    -Description $description `
    -BinaryPathName "$binaryPath -config=$configPath -service" `
    -StartupType Automatic

Write-Host "Service created successfully!"
Write-Host "To start the service, run: Start-Service -Name $serviceName"
