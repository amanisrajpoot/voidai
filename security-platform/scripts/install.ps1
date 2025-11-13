# Universal Security Platform Installation Script for Windows
# PowerShell script for Windows installation

$ErrorActionPreference = 'Stop'

$Version = "1.0.0"
$ControlPlaneUrl = $env:SECURITY_PLATFORM_CONTROL_PLANE_URL
if (-not $ControlPlaneUrl) {
    $ControlPlaneUrl = "https://api.securityplatform.com"
}
$AuthKey = $env:SECURITY_PLATFORM_AUTH_KEY

Write-Host "==========================================" -ForegroundColor Green
Write-Host "Security Platform Agent Installer" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Green
Write-Host ""

# Check for Chocolatey
if (Get-Command choco -ErrorAction SilentlyContinue) {
    Write-Host "[INFO] Installing via Chocolatey..." -ForegroundColor Green
    choco install security-platform-agent -y
    exit 0
}

# Fallback to direct MSI download
Write-Host "[INFO] Downloading installer..." -ForegroundColor Green
$DownloadUrl = "https://releases.securityplatform.com/security-platform-agent-${Version}.msi"
$TempFile = Join-Path $env:TEMP "security-platform-agent.msi"

try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempFile -UseBasicParsing
    
    Write-Host "[INFO] Installing..." -ForegroundColor Green
    Start-Process msiexec.exe -Wait -ArgumentList "/i `"$TempFile`" /quiet /norestart"
    
    # Create config directory
    $ConfigDir = Join-Path $env:ProgramData "SecurityPlatform"
    if (-not (Test-Path $ConfigDir)) {
        New-Item -ItemType Directory -Path $ConfigDir | Out-Null
    }
    
    # Create config file
    $ConfigFile = Join-Path $ConfigDir "agent-config.yaml"
    if (-not (Test-Path $ConfigFile)) {
        @"
control_plane_url: "$ControlPlaneUrl"
auth_key: "$AuthKey"
service_name: "desktop-agent"
environment: "production"
"@ | Out-File -FilePath $ConfigFile -Encoding UTF8
    }
    
    # Start service
    $ServiceName = "SecurityPlatformAgent"
    $Service = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
    if ($Service) {
        if ($Service.Status -ne 'Running') {
            Start-Service -Name $ServiceName
        }
    }
    
    Write-Host ""
    Write-Host "==========================================" -ForegroundColor Green
    Write-Host "Installation complete!" -ForegroundColor Green
    Write-Host "==========================================" -ForegroundColor Green
    Write-Host "Configuration: $ConfigFile" -ForegroundColor Yellow
    Write-Host ""
    
} catch {
    Write-Host "[ERROR] Installation failed: $_" -ForegroundColor Red
    exit 1
} finally {
    if (Test-Path $TempFile) {
        Remove-Item $TempFile -Force
    }
}
