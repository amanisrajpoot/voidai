# PowerShell Universal Installer for Security Platform Agent

$ErrorActionPreference = 'Stop'

Write-Host "Security Platform Agent Installer v1.0.0" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Green
Write-Host ""

# Check if running as administrator
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if (-not $isAdmin) {
    Write-Host "This script requires administrator privileges. Please run as administrator." -ForegroundColor Red
    exit 1
}

# Check for Chocolatey
if (Get-Command choco -ErrorAction SilentlyContinue) {
    Write-Host "Installing via Chocolatey..." -ForegroundColor Cyan
    choco install security-platform-agent -y
} else {
    Write-Host "Chocolatey not found. Installing via direct download..." -ForegroundColor Yellow
    
    $downloadUrl = "https://releases.securityplatform.com/agent.msi"
    $installerPath = "$env:TEMP\security-platform-agent.msi"
    
    Write-Host "Downloading installer..." -ForegroundColor Cyan
    Invoke-WebRequest -Uri $downloadUrl -OutFile $installerPath
    
    Write-Host "Installing..." -ForegroundColor Cyan
    Start-Process msiexec.exe -Wait -ArgumentList "/i $installerPath /quiet /norestart"
    
    Remove-Item $installerPath -Force
}

# Configuration
$controlPlaneUrl = $env:SECURITY_PLATFORM_CONTROL_PLANE_URL
$authKey = $env:SECURITY_PLATFORM_AUTH_KEY
$serviceName = $env:SECURITY_PLATFORM_SERVICE_NAME

if (-not $controlPlaneUrl -or -not $authKey) {
    Write-Host ""
    Write-Host "Configuration required:" -ForegroundColor Yellow
    $controlPlaneUrl = Read-Host "Control Plane URL"
    $authKey = Read-Host "Auth Key"
    $serviceName = Read-Host "Service Name (optional)"
}

$configPath = "$env:ProgramData\SecurityPlatform"
New-Item -ItemType Directory -Force -Path $configPath | Out-Null

$configFile = Join-Path $configPath "config.yaml"
@"
control_plane_url: $controlPlaneUrl
auth_key: $authKey
service_name: $($serviceName ?? 'agent')
environment: production
otlp_endpoint: http://localhost:4318/v1/traces
local_policy: observe
"@ | Out-File -FilePath $configFile -Encoding UTF8

# Restart service
$serviceName = "SecurityPlatformAgent"
Restart-Service -Name $serviceName -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "Installation complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "1. Verify service: Get-Service SecurityPlatformAgent"
Write-Host "2. View logs: Get-EventLog -LogName Application -Source SecurityPlatformAgent"
Write-Host "3. Edit config: notepad $configFile"
