$ErrorActionPreference = 'Stop'

$serviceName = 'SecurityPlatformAgent'
$installDir = Join-Path $env:ProgramFiles "SecurityPlatform"
$exePath = Join-Path $installDir "security-platform-agent.exe"

# Stop and remove service
if (Get-Service -Name $serviceName -ErrorAction SilentlyContinue) {
  Stop-Service -Name $serviceName -Force -ErrorAction SilentlyContinue
  & $exePath -remove
}

# Remove files
Remove-Item $installDir -Recurse -Force -ErrorAction SilentlyContinue

Write-Host "Security Platform Agent uninstalled successfully"
