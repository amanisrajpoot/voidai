$ErrorActionPreference = 'Stop'

$packageName = 'security-platform-agent'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$fileLocation = Join-Path $toolsDir 'security-platform-agent.exe'

$packageArgs = @{
  packageName   = $packageName
  fileType      = 'EXE'
  file          = $fileLocation
  silentArgs    = '/S'
  validExitCodes= @(0)
}

Install-ChocolateyPackage @packageArgs

# Install service
$serviceName = 'SecurityPlatformAgent'
$serviceExists = Get-Service -Name $serviceName -ErrorAction SilentlyContinue

if (-not $serviceExists) {
    $exePath = Join-Path $env:ProgramFiles 'SecurityPlatform\security-platform-agent.exe'
    New-Service -Name $serviceName `
                -BinaryPathName $exePath `
                -DisplayName 'Security Platform Agent' `
                -Description 'Security Platform Desktop Agent for Windows' `
                -StartupType Automatic
}

# Create config directory
$configDir = Join-Path $env:ProgramData 'SecurityPlatform'
if (-not (Test-Path $configDir)) {
    New-Item -ItemType Directory -Path $configDir | Out-Null
}

$configFile = Join-Path $configDir 'agent-config.yaml'
if (-not (Test-Path $configFile)) {
    @"
control_plane_url: "https://api.securityplatform.com"
auth_key: ""
service_name: "desktop-agent"
environment: "production"
"@ | Out-File -FilePath $configFile -Encoding UTF8
}
