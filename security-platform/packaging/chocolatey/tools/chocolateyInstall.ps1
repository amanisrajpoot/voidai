$ErrorActionPreference = 'Stop'

$packageName = 'security-platform-agent'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$installPath = "$env:ProgramFiles\SecurityPlatform\Agent"

# Create installation directory
New-Item -ItemType Directory -Force -Path $installPath | Out-Null

# Install binary
$file = Join-Path $toolsDir "security-platform-agent.exe"
Copy-Item $file -Destination "$installPath\security-platform-agent.exe" -Force

# Create config directory
$configPath = "$env:ProgramData\SecurityPlatform"
New-Item -ItemType Directory -Force -Path $configPath | Out-Null

# Create default config if it doesn't exist
$configFile = Join-Path $configPath "config.yaml"
if (-not (Test-Path $configFile)) {
    @"
service_name: desktop-agent
version: 1.0.0
environment: production
otlp_endpoint: http://localhost:4318/v1/traces
local_policy: observe
"@ | Out-File -FilePath $configFile -Encoding UTF8
}

# Install as Windows Service
$serviceName = "SecurityPlatformAgent"
$existingService = Get-Service -Name $serviceName -ErrorAction SilentlyContinue

if ($existingService) {
    Stop-Service -Name $serviceName -Force -ErrorAction SilentlyContinue
    sc.exe delete $serviceName | Out-Null
    Start-Sleep -Seconds 2
}

New-Service -Name $serviceName `
    -DisplayName "Security Platform Desktop Agent" `
    -Description "Security Platform Desktop Agent for observability and security" `
    -BinaryPathName "$installPath\security-platform-agent.exe -service" `
    -StartupType Automatic | Out-Null

Start-Service -Name $serviceName

Write-Host "Security Platform Agent installed successfully!" -ForegroundColor Green
