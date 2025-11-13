package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var packageTargets string
var packageOut string

// packageCmd represents the package command
var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Produce installers for specified targets",
	Long: `Generate installers for the specified package formats.
Supported targets: deb, rpm, msi, dmg, helm, airgap`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if packageTargets == "" {
			return fmt.Errorf("--target is required")
		}

		targets := strings.Split(packageTargets, ",")
		outDir := packageOut
		if outDir == "" {
			outDir = viper.GetString("out")
		}

		if err := os.MkdirAll(outDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		for _, target := range targets {
			target = strings.TrimSpace(target)
			fmt.Printf("Packaging for %s...\n", target)

			switch target {
			case "deb":
				if err := packageDeb(outDir); err != nil {
					return fmt.Errorf("failed to package deb: %w", err)
				}
			case "rpm":
				if err := packageRpm(outDir); err != nil {
					return fmt.Errorf("failed to package rpm: %w", err)
				}
			case "msi":
				if err := packageMsi(outDir); err != nil {
					return fmt.Errorf("failed to package msi: %w", err)
				}
			case "dmg":
				if err := packageDmg(outDir); err != nil {
					return fmt.Errorf("failed to package dmg: %w", err)
				}
			case "pkg":
				if err := packagePkg(outDir); err != nil {
					return fmt.Errorf("failed to package pkg: %w", err)
				}
			case "helm":
				if err := packageHelm(outDir); err != nil {
					return fmt.Errorf("failed to package helm: %w", err)
				}
			case "homebrew":
				if err := packageHomebrew(outDir); err != nil {
					return fmt.Errorf("failed to package homebrew: %w", err)
				}
			case "chocolatey":
				if err := packageChocolatey(outDir); err != nil {
					return fmt.Errorf("failed to package chocolatey: %w", err)
				}
			case "airgap":
				if err := packageAirgap(outDir); err != nil {
					return fmt.Errorf("failed to package airgap: %w", err)
				}
			default:
				return fmt.Errorf("unsupported target: %s", target)
			}
		}

		fmt.Printf("Packages created in %s\n", outDir)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(packageCmd)
	packageCmd.Flags().StringVar(&packageTargets, "target", "", "comma-separated list of targets (deb,rpm,msi,dmg,pkg,helm,homebrew,chocolatey,airgap)")
	packageCmd.Flags().StringVar(&packageOut, "out", "", "output directory (default: ./dist)")
}

func packageDeb(outDir string) error {
	// Create DEB package structure
	debDir := filepath.Join(outDir, "deb")
	if err := os.MkdirAll(filepath.Join(debDir, "DEBIAN"), 0755); err != nil {
		return err
	}

	// Write control file
	control := `Package: security-platform-agent
Version: 1.0.0
Architecture: amd64
Maintainer: Security Platform <support@securityplatform.com>
Description: Security Platform Agent
`
	if err := os.WriteFile(filepath.Join(debDir, "DEBIAN", "control"), []byte(control), 0644); err != nil {
		return err
	}

	fmt.Println("  Created DEB package structure")
	return nil
}

func packageRpm(outDir string) error {
	// Create RPM spec file
	spec := `Name: security-platform-agent
Version: 1.0.0
Release: 1
Summary: Security Platform Agent
License: MIT

%description
Security Platform Agent for observability and security.

%files
`
	specPath := filepath.Join(outDir, "security-platform-agent.spec")
	if err := os.WriteFile(specPath, []byte(spec), 0644); err != nil {
		return err
	}

	fmt.Println("  Created RPM spec file")
	return nil
}

func packageMsi(outDir string) error {
	// Create WiX XML for MSI
	wix := `<?xml version="1.0" encoding="UTF-8"?>
<Wix xmlns="http://schemas.microsoft.com/wix/2006/wi">
  <Product Id="*" Name="Security Platform Agent" Language="1033" Version="1.0.0" Manufacturer="Security Platform">
    <Package InstallerVersion="200" Compressed="yes" />
    <MediaTemplate />
    <Feature Id="ProductFeature" Title="Security Platform Agent" Level="1">
      <ComponentRef Id="AgentBinary" />
    </Feature>
  </Product>
</Wix>
`
	wixPath := filepath.Join(outDir, "agent.wxs")
	if err := os.WriteFile(wixPath, []byte(wix), 0644); err != nil {
		return err
	}

	fmt.Println("  Created WiX XML file")
	return nil
}

func packageDmg(outDir string) error {
	// Create DMG build script
	script := `#!/bin/bash
# DMG packaging script
hdiutil create -volname "Security Platform Agent" -srcfolder ./build -ov -format UDZO security-platform-agent.dmg
`
	scriptPath := filepath.Join(outDir, "build-dmg.sh")
	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		return err
	}

	fmt.Println("  Created DMG build script")
	return nil
}

func packagePkg(outDir string) error {
	// Create PKG build script for macOS
	script := `#!/bin/bash
# PKG packaging script for macOS
pkgbuild --root ./build --identifier com.securityplatform.agent --version 1.0.0 --install-location /usr/local security-platform-agent.pkg
`
	scriptPath := filepath.Join(outDir, "build-pkg.sh")
	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		return err
	}

	fmt.Println("  Created PKG build script")
	return nil
}

func packageHomebrew(outDir string) error {
	homebrewDir := filepath.Join(outDir, "homebrew")
	if err := os.MkdirAll(homebrewDir, 0755); err != nil {
		return err
	}

	// Copy Homebrew formula
	formulaPath := filepath.Join("..", "..", "packaging", "homebrew", "security-platform-agent.rb")
	destPath := filepath.Join(homebrewDir, "security-platform-agent.rb")
	
	// Read and update formula
	formula := `class SecurityPlatformAgent < Formula
  desc "Security Platform Desktop Agent"
  homepage "https://github.com/security-platform/desktop-agent"
  url "https://github.com/security-platform/desktop-agent/releases/download/v1.0.0/security-platform-agent-darwin-amd64.tar.gz"
  sha256 "placeholder-sha256"
  version "1.0.0"

  def install
    bin.install "security-platform-agent"
  end

  service do
    run [opt_bin/"security-platform-agent", "-service"]
    keep_alive true
  end
end
`
	if err := os.WriteFile(destPath, []byte(formula), 0644); err != nil {
		return err
	}

	fmt.Println("  Created Homebrew formula")
	return nil
}

func packageChocolatey(outDir string) error {
	chocoDir := filepath.Join(outDir, "chocolatey")
	if err := os.MkdirAll(chocoDir, 0755); err != nil {
		return err
	}

	// Create Chocolatey package structure
	nuspec := `<?xml version="1.0" encoding="utf-8"?>
<package xmlns="http://schemas.microsoft.com/packaging/2015/06/nuspec.xsd">
  <metadata>
    <id>security-platform-agent</id>
    <version>1.0.0</version>
    <title>Security Platform Agent</title>
    <authors>Security Platform</authors>
    <summary>Security Platform Desktop Agent for Windows</summary>
    <description>Simple, automatic security observability for Windows desktop systems.</description>
    <projectUrl>https://github.com/security-platform/desktop-agent</projectUrl>
    <tags>security observability monitoring agent</tags>
    <requireLicenseAcceptance>false</requireLicenseAcceptance>
  </metadata>
  <files>
    <file src="tools\**" target="tools" />
  </files>
</package>
`
	nuspecPath := filepath.Join(chocoDir, "security-platform-agent.nuspec")
	if err := os.WriteFile(nuspecPath, []byte(nuspec), 0644); err != nil {
		return err
	}

	// Create tools directory
	toolsDir := filepath.Join(chocoDir, "tools")
	if err := os.MkdirAll(toolsDir, 0755); err != nil {
		return err
	}

	fmt.Println("  Created Chocolatey package")
	return nil
}

func packageHelm(outDir string) error {
	helmDir := filepath.Join(outDir, "helm", "security-platform-agent")
	if err := os.MkdirAll(helmDir, 0755); err != nil {
		return err
	}

	// Create Chart.yaml
	chart := `apiVersion: v2
name: security-platform-agent
description: Security Platform Agent Helm Chart
type: application
version: 1.0.0
appVersion: "1.0.0"
`
	if err := os.WriteFile(filepath.Join(helmDir, "Chart.yaml"), []byte(chart), 0644); err != nil {
		return err
	}

	// Create values.yaml
	values := `replicaCount: 1

image:
  repository: security-platform/agent
  tag: "1.0.0"

config:
  controlPlaneUrl: ""
  otlpEndpoint: "http://localhost:4318"
`
	if err := os.WriteFile(filepath.Join(helmDir, "values.yaml"), []byte(values), 0644); err != nil {
		return err
	}

	fmt.Println("  Created Helm chart")
	return nil
}

func packageAirgap(outDir string) error {
	airgapDir := filepath.Join(outDir, "airgap")
	if err := os.MkdirAll(airgapDir, 0755); err != nil {
		return err
	}

	// Create comprehensive installation script
	installScript := `#!/bin/bash
# Air-gapped installation script for Security Platform
set -e

echo "=========================================="
echo "Security Platform Air-Gap Installation"
echo "=========================================="

# Detect OS
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$ID
    VERSION=$VERSION_ID
elif [ -f /etc/debian_version ]; then
    OS="debian"
elif [ -f /etc/redhat-release ]; then
    OS="rhel"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    OS="darwin"
else
    echo "Unsupported OS"
    exit 1
fi

echo "Detected OS: $OS"

# Extract packages
if [ -f packages.tar.gz ]; then
    echo "Extracting packages..."
    tar -xzf packages.tar.gz
fi

# Install based on OS
case $OS in
    debian|ubuntu)
        echo "Installing DEB package..."
        if [ -f security-platform-agent*.deb ]; then
            sudo dpkg -i security-platform-agent*.deb || sudo apt-get install -f -y
        fi
        ;;
    rhel|centos|fedora|rocky|almalinux)
        echo "Installing RPM package..."
        if [ -f security-platform-agent*.rpm ]; then
            sudo rpm -ivh security-platform-agent*.rpm || sudo yum install -y security-platform-agent*.rpm
        fi
        ;;
    darwin)
        echo "Installing macOS package..."
        if [ -f security-platform-agent*.pkg ]; then
            sudo installer -pkg security-platform-agent*.pkg -target /
        fi
        ;;
    *)
        echo "Manual installation required for $OS"
        exit 1
        ;;
esac

# Create config directory
CONFIG_DIR="/etc/security-platform"
if [ ! -d "$CONFIG_DIR" ]; then
    echo "Creating config directory..."
    sudo mkdir -p "$CONFIG_DIR"
fi

# Copy default config if it doesn't exist
if [ -f agent-config.yaml ] && [ ! -f "$CONFIG_DIR/agent-config.yaml" ]; then
    echo "Installing default configuration..."
    sudo cp agent-config.yaml "$CONFIG_DIR/"
fi

# Enable and start service (Linux)
if [ "$OS" != "darwin" ]; then
    echo "Enabling service..."
    sudo systemctl enable security-platform-agent || true
    echo "Starting service..."
    sudo systemctl start security-platform-agent || true
fi

echo ""
echo "=========================================="
echo "Installation complete!"
echo "=========================================="
echo "Configuration: $CONFIG_DIR/agent-config.yaml"
echo "Logs: /var/log/security-platform-agent.log"
echo ""
`
	if err := os.WriteFile(filepath.Join(airgapDir, "install.sh"), []byte(installScript), 0755); err != nil {
		return err
	}

	// Create comprehensive manifest
	manifest := `{
  "version": "1.0.0",
  "components": {
    "agents": {
      "java": "1.0.0",
      "dotnet": "1.0.0",
      "go": "1.0.0",
      "node": "1.0.0",
      "python": "1.0.0"
    },
    "sdks": {
      "ios": "1.0.0",
      "android": "1.0.0",
      "frontend": "1.0.0"
    },
    "desktop": {
      "linux": "1.0.0",
      "darwin": "1.0.0",
      "windows": "1.0.0"
    },
    "operator": "1.0.0"
  },
  "packages": {
    "linux": [
      "security-platform-agent.deb",
      "security-platform-agent.rpm"
    ],
    "darwin": [
      "security-platform-agent.pkg"
    ],
    "windows": [
      "security-platform-agent.msi"
    ]
  },
  "dependencies": [],
  "checksums": {},
  "installation_instructions": "Run ./install.sh as root or with sudo"
}
`
	if err := os.WriteFile(filepath.Join(airgapDir, "manifest.json"), []byte(manifest), 0644); err != nil {
		return err
	}

	// Create README for airgap bundle
	readme := `# Security Platform Air-Gap Bundle

This bundle contains all components needed for air-gapped installation.

## Contents

- Agent packages for all platforms (Linux, macOS, Windows)
- SDKs for all languages (Java, .NET, Go, Node.js, Python)
- Mobile SDKs (iOS, Android)
- Kubernetes Operator
- Installation scripts

## Installation

1. Transfer this bundle to your air-gapped system
2. Extract: \`tar -xzf airgap-bundle.tar.gz\`
3. Run installer: \`sudo ./install.sh\`
4. Configure: Edit \`/etc/security-platform/agent-config.yaml\`
5. Start service: \`sudo systemctl start security-platform-agent\`

## Manual Installation

If automatic installation fails, refer to platform-specific documentation:
- Linux: See DEB/RPM package installation
- macOS: See PKG installation
- Windows: See MSI installation

## Support

For air-gap installation support, contact: support@securityplatform.com
`
	if err := os.WriteFile(filepath.Join(airgapDir, "README.md"), []byte(readme), 0644); err != nil {
		return err
	}

	fmt.Println("  Created airgap bundle")
	return nil
}
