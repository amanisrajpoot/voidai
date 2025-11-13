# Chocolatey Package for Security Platform Agent

## Installation

```powershell
# Install from local package
choco install security-platform-agent -s . --version 1.0.0

# Or from Chocolatey repository (once published)
choco install security-platform-agent
```

## Configuration

Edit the configuration file:

```powershell
notepad C:\ProgramData\SecurityPlatform\config.yaml
```

## Usage

```powershell
# Start service
Start-Service -Name SecurityPlatformAgent

# Stop service
Stop-Service -Name SecurityPlatformAgent

# Check status
Get-Service -Name SecurityPlatformAgent
```

## Updating

```powershell
choco upgrade security-platform-agent
```

## Uninstallation

```powershell
choco uninstall security-platform-agent
```
