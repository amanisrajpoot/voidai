$ErrorActionPreference = 'Stop'

$serviceName = "SecurityPlatformAgent"
$installPath = "$env:ProgramFiles\SecurityPlatform\Agent"

# Stop and remove service
$service = Get-Service -Name $serviceName -ErrorAction SilentlyContinue
if ($service) {
    Stop-Service -Name $serviceName -Force -ErrorAction SilentlyContinue
    sc.exe delete $serviceName | Out-Null
}

# Remove installation directory
if (Test-Path $installPath) {
    Remove-Item -Recurse -Force $installPath
}

Write-Host "Security Platform Agent uninstalled successfully!" -ForegroundColor Green
