# Installation Scripts

Universal installation scripts for all platforms.

## Quick Install

### Linux/macOS

```bash
curl -fsSL https://install.securityplatform.com | bash
```

Or with environment variables:

```bash
export SECURITY_PLATFORM_CONTROL_PLANE_URL=https://api.securityplatform.com
export SECURITY_PLATFORM_AUTH_KEY=your-token-here
curl -fsSL https://install.securityplatform.com | bash
```

### Windows (PowerShell)

```powershell
$env:SECURITY_PLATFORM_CONTROL_PLANE_URL="https://api.securityplatform.com"
$env:SECURITY_PLATFORM_AUTH_KEY="your-token-here"
iex (New-Object Net.WebClient).DownloadString('https://install.securityplatform.com/install.ps1')
```

## Manual Installation

### Linux

```bash
wget https://install.securityplatform.com/install.sh
chmod +x install.sh
sudo ./install.sh
```

### Windows

```powershell
Invoke-WebRequest -Uri https://install.securityplatform.com/install.ps1 -OutFile install.ps1
.\install.ps1
```

## Features

- ✅ Automatic OS detection
- ✅ Package manager integration (Homebrew, Chocolatey)
- ✅ Automatic service setup
- ✅ Configuration file creation
- ✅ One-command installation
