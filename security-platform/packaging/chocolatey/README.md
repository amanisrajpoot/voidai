# Chocolatey Package

Chocolatey package for Windows installation.

## Building

```powershell
choco pack security-platform-agent.nuspec
```

## Installation

```powershell
choco install security-platform-agent -y
```

## Configuration

Edit configuration at:
`C:\ProgramData\SecurityPlatform\agent-config.yaml`

## Service Management

```powershell
# Start service
net start SecurityPlatformAgent

# Stop service
net stop SecurityPlatformAgent

# Check status
sc query SecurityPlatformAgent
```
