# Homebrew Formula for Security Platform Agent

## Installation

```bash
# Tap the repository
brew tap securityplatform/agent

# Install the agent
brew install security-platform-agent

# Start the service
brew services start security-platform-agent
```

## Configuration

Edit the configuration file:

```bash
sudo nano /usr/local/etc/security-platform/config.yaml
```

## Usage

```bash
# Start service
brew services start security-platform-agent

# Stop service
brew services stop security-platform-agent

# Restart service
brew services restart security-platform-agent

# Check status
brew services list | grep security-platform-agent
```

## Updating

```bash
brew update
brew upgrade security-platform-agent
```
