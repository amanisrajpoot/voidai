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
			case "helm":
				if err := packageHelm(outDir); err != nil {
					return fmt.Errorf("failed to package helm: %w", err)
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
	packageCmd.Flags().StringVar(&packageTargets, "target", "", "comma-separated list of targets (deb,rpm,msi,dmg,helm,airgap)")
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

	// Create installation script
	installScript := `#!/bin/bash
# Air-gapped installation script
set -e

echo "Installing Security Platform Agent in air-gapped environment..."

# Extract packages
tar -xzf packages.tar.gz

# Install based on OS
if [ -f /etc/debian_version ]; then
    dpkg -i security-platform-agent*.deb
elif [ -f /etc/redhat-release ]; then
    rpm -ivh security-platform-agent*.rpm
fi

echo "Installation complete."
`
	if err := os.WriteFile(filepath.Join(airgapDir, "install.sh"), []byte(installScript), 0755); err != nil {
		return err
	}

	// Create manifest
	manifest := `{
  "version": "1.0.0",
  "packages": [
    "security-platform-agent.deb",
    "security-platform-agent.rpm"
  ],
  "dependencies": [],
  "checksums": {}
}
`
	if err := os.WriteFile(filepath.Join(airgapDir, "manifest.json"), []byte(manifest), 0644); err != nil {
		return err
	}

	fmt.Println("  Created airgap bundle")
	return nil
}
