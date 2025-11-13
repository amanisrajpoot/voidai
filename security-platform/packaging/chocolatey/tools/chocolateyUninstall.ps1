$ErrorActionPreference = 'Stop'

$serviceName = 'SecurityPlatformAgent'
$service = Get-Service -Name $serviceName -ErrorAction SilentlyContinue

if ($service) {
    if ($service.Status -eq 'Running') {
        Stop-Service -Name $serviceName -Force
    }
    Remove-Service -Name $serviceName
}

$configDir = Join-Path $env:ProgramData 'SecurityPlatform'
if (Test-Path $configDir) {
    Remove-Item -Path $configDir -Recurse -Force
}
