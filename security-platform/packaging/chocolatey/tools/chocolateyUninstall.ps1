$ErrorActionPreference = 'Stop'

$serviceName = "SecurityPlatformAgent"
$installDir = "$env:ProgramFiles\SecurityPlatform"

# Stop and remove service
$service = Get-Service -Name $serviceName -ErrorAction SilentlyContinue
if ($service) {
    Write-Host "Stopping service..."
    Stop-Service -Name $serviceName -Force -ErrorAction SilentlyContinue
    sc.exe delete $serviceName
}

# Remove installation directory
if (Test-Path $installDir) {
    Remove-Item -Recurse -Force $installDir
}

Write-Host "Uninstallation complete!"
