# Chocolatey Package

## Building

```powershell
choco pack security-platform-agent.nuspec
```

## Installing

```powershell
choco install security-platform-agent -s .
```

## Publishing

```powershell
choco push security-platform-agent.1.0.0.nupkg --source https://push.chocolatey.org/
```
