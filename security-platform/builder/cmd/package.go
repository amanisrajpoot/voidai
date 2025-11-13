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
Supported targets: deb, rpm, msi, dmg, helm, homebrew, chocolatey, kubernetes, airgap`,
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
			case "helm":
				if err := packageHelm(outDir); err != nil {
					return fmt.Errorf("failed to package helm: %w", err)
				}
			case "airgap":
				if err := packageAirgap(outDir); err != nil {
					return fmt.Errorf("failed to package airgap: %w", err)
				}
			case "homebrew":
				if err := packageHomebrew(outDir); err != nil {
					return fmt.Errorf("failed to package homebrew: %w", err)
				}
			case "chocolatey":
				if err := packageChocolatey(outDir); err != nil {
					return fmt.Errorf("failed to package chocolatey: %w", err)
				}
			case "kubernetes":
				if err := packageKubernetes(outDir); err != nil {
					return fmt.Errorf("failed to package kubernetes: %w", err)
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
	packageCmd.Flags().StringVar(&packageTargets, "target", "", "comma-separated list of targets (deb,rpm,msi,dmg,helm,homebrew,chocolatey,kubernetes,airgap)")
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

	// Create packages directory
	packagesDir := filepath.Join(airgapDir, "packages")
	if err := os.MkdirAll(packagesDir, 0755); err != nil {
		return err
	}

	// Create installation script for Linux
	installScript := `#!/bin/bash
# Air-gapped installation script
set -e

echo "Installing Security Platform Agent in air-gapped environment..."

# Detect OS
if [ -f /etc/debian_version ]; then
    echo "Detected Debian/Ubuntu system"
    if [ -f packages/security-platform-agent*.deb ]; then
        dpkg -i packages/security-platform-agent*.deb
    else
        echo "ERROR: No .deb package found"
        exit 1
    fi
elif [ -f /etc/redhat-release ]; then
    echo "Detected RedHat/CentOS system"
    if [ -f packages/security-platform-agent*.rpm ]; then
        rpm -ivh packages/security-platform-agent*.rpm
    else
        echo "ERROR: No .rpm package found"
        exit 1
    fi
elif [ "$(uname)" == "Darwin" ]; then
    echo "Detected macOS system"
    if [ -f packages/security-platform-agent*.pkg ]; then
        sudo installer -pkg packages/security-platform-agent*.pkg -target /
    elif [ -f packages/security-platform-agent ]; then
        sudo cp packages/security-platform-agent /usr/local/bin/
        sudo chmod +x /usr/local/bin/security-platform-agent
    else
        echo "ERROR: No macOS package found"
        exit 1
    fi
else
    echo "ERROR: Unsupported operating system"
    exit 1
fi

# Create config directory
sudo mkdir -p /etc/security-platform

# Copy default config if it doesn't exist
if [ ! -f /etc/security-platform/agent-config.yaml ]; then
    sudo cp config/agent-config.yaml /etc/security-platform/
fi

echo "Installation complete."
echo "Configuration: /etc/security-platform/agent-config.yaml"
`
	if err := os.WriteFile(filepath.Join(airgapDir, "install.sh"), []byte(installScript), 0755); err != nil {
		return err
	}

	// Create Windows installation script
	installWindowsScript := `@echo off
REM Air-gapped installation script for Windows
echo Installing Security Platform Agent in air-gapped environment...

if exist packages\security-platform-agent*.msi (
    echo Installing MSI package...
    msiexec /i packages\security-platform-agent*.msi /quiet /norestart
) else if exist packages\security-platform-agent.exe (
    echo Installing executable...
    copy packages\security-platform-agent.exe "%ProgramFiles%\SecurityPlatform\"
    "%ProgramFiles%\SecurityPlatform\security-platform-agent.exe" -install
) else (
    echo ERROR: No Windows package found
    exit /b 1
)

REM Create config directory
if not exist "%ProgramData%\SecurityPlatform" mkdir "%ProgramData%\SecurityPlatform"

REM Copy default config if it doesn't exist
if not exist "%ProgramData%\SecurityPlatform\agent-config.yaml" (
    copy config\agent-config.yaml "%ProgramData%\SecurityPlatform\"
)

echo Installation complete.
echo Configuration: %ProgramData%\SecurityPlatform\agent-config.yaml
`
	if err := os.WriteFile(filepath.Join(airgapDir, "install.bat"), []byte(installWindowsScript), 0644); err != nil {
		return err
	}

	// Create manifest with all package types
	manifest := `{
  "version": "1.0.0",
  "packages": {
    "linux": {
      "deb": "security-platform-agent.deb",
      "rpm": "security-platform-agent.rpm"
    },
    "darwin": {
      "pkg": "security-platform-agent.pkg",
      "binary": "security-platform-agent"
    },
    "windows": {
      "msi": "security-platform-agent.msi",
      "exe": "security-platform-agent.exe"
    }
  },
  "dependencies": [],
  "checksums": {},
  "installation": {
    "linux": "./install.sh",
    "darwin": "./install.sh",
    "windows": "install.bat"
  }
}
`
	if err := os.WriteFile(filepath.Join(airgapDir, "manifest.json"), []byte(manifest), 0644); err != nil {
		return err
	}

	// Create README
	readme := `# Air-Gapped Installation Bundle

This bundle contains all necessary files for installing Security Platform Agent in an air-gapped environment.

## Contents

- \`packages/\` - Platform-specific installation packages
- \`config/\` - Default configuration files
- \`install.sh\` - Linux/macOS installation script
- \`install.bat\` - Windows installation script
- \`manifest.json\` - Bundle manifest

## Installation

### Linux

\`\`\`bash
chmod +x install.sh
sudo ./install.sh
\`\`\`

### macOS

\`\`\`bash
chmod +x install.sh
sudo ./install.sh
\`\`\`

### Windows

Run \`install.bat\` as Administrator.

## Configuration

After installation, edit the configuration file:

- Linux/macOS: \`/etc/security-platform/agent-config.yaml\`
- Windows: \`C:\\ProgramData\\SecurityPlatform\\agent-config.yaml\`

## Verification

\`\`\`bash
# Linux/macOS
security-platform-agent --version

# Windows
security-platform-agent.exe --version
\`\`\`
`
	if err := os.WriteFile(filepath.Join(airgapDir, "README.md"), []byte(readme), 0644); err != nil {
		return err
	}

	// Create config directory
	configDir := filepath.Join(airgapDir, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Copy default config template
	defaultConfig := `control_plane_url: "https://api.securityplatform.com"
auth_key: ""  # Set this after installation
service_name: "security-platform-agent"
environment: "production"
namespace: "default"

telemetry:
  otlp_endpoint: "http://localhost:4318"
  batch_size: 100
  batch_timeout: "5s"
  export_timeout: "30s"
  max_queue_size: 2048

policy:
  mode: "observe"
  auto_enable_blocking: false
  observe_period_hours: 48

security:
  mtls_enabled: true
`
	if err := os.WriteFile(filepath.Join(configDir, "agent-config.yaml"), []byte(defaultConfig), 0644); err != nil {
		return err
	}

	fmt.Println("  Created airgap bundle with installation scripts")
	return nil
}

func packageHomebrew(outDir string) error {
	homebrewDir := filepath.Join(outDir, "homebrew")
	if err := os.MkdirAll(homebrewDir, 0755); err != nil {
		return err
	}

	// Copy Homebrew formula
	formulaPath := filepath.Join("packaging", "homebrew", "security-platform-agent.rb")
	if _, err := os.Stat(formulaPath); err == nil {
		if err := copyFile(formulaPath, filepath.Join(homebrewDir, "security-platform-agent.rb")); err != nil {
			return err
		}
	}

	fmt.Println("  Created Homebrew formula")
	return nil
}

func packageChocolatey(outDir string) error {
	chocoDir := filepath.Join(outDir, "chocolatey")
	if err := os.MkdirAll(chocoDir, 0755); err != nil {
		return err
	}

	// Copy Chocolatey package files
	nuspecPath := filepath.Join("packaging", "chocolatey", "security-platform-agent.nuspec")
	if _, err := os.Stat(nuspecPath); err == nil {
		if err := copyFile(nuspecPath, filepath.Join(chocoDir, "security-platform-agent.nuspec")); err != nil {
			return err
		}
	}

	toolsDir := filepath.Join(chocoDir, "tools")
	if err := os.MkdirAll(toolsDir, 0755); err != nil {
		return err
	}

	installScript := filepath.Join("packaging", "chocolatey", "tools", "chocolateyInstall.ps1")
	if _, err := os.Stat(installScript); err == nil {
		if err := copyFile(installScript, filepath.Join(toolsDir, "chocolateyInstall.ps1")); err != nil {
			return err
		}
	}

	uninstallScript := filepath.Join("packaging", "chocolatey", "tools", "chocolateyUninstall.ps1")
	if _, err := os.Stat(uninstallScript); err == nil {
		if err := copyFile(uninstallScript, filepath.Join(toolsDir, "chocolateyUninstall.ps1")); err != nil {
			return err
		}
	}

	fmt.Println("  Created Chocolatey package")
	return nil
}

func packageKubernetes(outDir string) error {
	k8sDir := filepath.Join(outDir, "kubernetes")
	if err := os.MkdirAll(k8sDir, 0755); err != nil {
		return err
	}

	// Copy Kubernetes Operator files
	operatorDir := filepath.Join("packaging", "kubernetes-operator")
	if _, err := os.Stat(operatorDir); err == nil {
		// Copy CRDs
		crdDir := filepath.Join(operatorDir, "config", "crd", "bases")
		if _, err := os.Stat(crdDir); err == nil {
			if err := copyDir(crdDir, filepath.Join(k8sDir, "crd")); err != nil {
				return err
			}
		}

		// Copy examples
		examplesDir := filepath.Join(operatorDir, "examples")
		if _, err := os.Stat(examplesDir); err == nil {
			if err := copyDir(examplesDir, filepath.Join(k8sDir, "examples")); err != nil {
				return err
			}
		}
	}

	fmt.Println("  Created Kubernetes Operator package")
	return nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}

		return copyFile(path, dstPath)
	})
}
