$ErrorActionPreference = 'Stop'

$packageName = 'security-platform-agent'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$fileLocation = Join-Path $toolsDir "security-platform-agent.exe"

# Install binary
$installDir = "$env:ProgramFiles\SecurityPlatform"
New-Item -ItemType Directory -Force -Path $installDir | Out-Null
Copy-Item "$fileLocation" "$installDir\security-platform-agent.exe"

# Create config directory
$configDir = "$env:ProgramData\SecurityPlatform"
New-Item -ItemType Directory -Force -Path $configDir | Out-Null

# Copy example config if it doesn't exist
$configFile = "$configDir\config.yaml"
if (-not (Test-Path $configFile)) {
    Copy-Item "$toolsDir\config.yaml.example" $configFile
}

# Install Windows Service
$serviceName = "SecurityPlatformAgent"
$service = Get-Service -Name $serviceName -ErrorAction SilentlyContinue

if (-not $service) {
    Write-Host "Installing Windows Service..."
    $servicePath = "$installDir\security-platform-agent.exe"
    $serviceArgs = "-config=$configFile -service"
    
    New-Service -Name $serviceName `
        -DisplayName "Security Platform Agent" `
        -Description "Security observability agent for Windows" `
        -BinaryPathName "$servicePath $serviceArgs" `
        -StartupType Automatic | Out-Null
    
    Write-Host "Service installed successfully!"
} else {
    Write-Host "Service already exists. Updating..."
    Stop-Service -Name $serviceName -Force -ErrorAction SilentlyContinue
    sc.exe config $serviceName binPath= "$installDir\security-platform-agent.exe -config=$configFile -service"
}

Write-Host "Installation complete!"
Write-Host "Configuration file: $configFile"
Write-Host "To start the service: Start-Service -Name $serviceName"
