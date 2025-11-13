$ErrorActionPreference = 'Stop'

$packageName = 'security-platform-agent'
$url64 = 'https://github.com/security-platform/agents/releases/download/v1.0.0/security-platform-agent-1.0.0-windows-amd64.zip'
$checksum64 = 'placeholder-sha256'
$checksumType64 = 'sha256'

$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$installDir = Join-Path $env:ProgramFiles "SecurityPlatform"

# Create installation directory
New-Item -ItemType Directory -Force -Path $installDir | Out-Null

# Download and extract
$tempDir = Join-Path $env:TEMP "security-platform-agent"
New-Item -ItemType Directory -Force -Path $tempDir | Out-Null

$zipFile = Join-Path $tempDir "agent.zip"
Get-ChocolateyWebFile -PackageName $packageName `
  -FileFullPath $zipFile `
  -Url64bit $url64 `
  -Checksum64 $checksum64 `
  -ChecksumType64 $checksumType64

Expand-Archive -Path $zipFile -DestinationPath $tempDir -Force

# Copy files
Copy-Item "$tempDir\security-platform-agent.exe" -Destination "$installDir\security-platform-agent.exe" -Force

# Install as Windows Service
$exePath = Join-Path $installDir "security-platform-agent.exe"
& $exePath -install

# Create config directory
$configDir = Join-Path $env:ProgramData "SecurityPlatform"
New-Item -ItemType Directory -Force -Path $configDir | Out-Null

# Create default config
$configFile = Join-Path $configDir "agent-config.yaml"
if (-not (Test-Path $configFile)) {
  @"
control_plane_url: "https://api.securityplatform.com"
service_name: "security-platform-agent"
environment: "production"
"@ | Out-File -FilePath $configFile -Encoding UTF8
}

# Cleanup
Remove-Item $tempDir -Recurse -Force

Write-Host "Security Platform Agent installed successfully"
Write-Host "Configuration file: $configFile"
Write-Host "To start the service: net start SecurityPlatformAgent"
