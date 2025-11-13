package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var airgapOutput string
var airgapVersion string
var airgapIncludeAgents []string
var airgapIncludePlatforms []string

// airgapCmd represents the airgap command
var airgapCmd = &cobra.Command{
	Use:   "airgap",
	Short: "Generate airgap installation bundle",
	Long: `Generate a complete offline installation bundle containing all agents,
SDKs, installers, and dependencies for air-gapped environments.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if airgapVersion == "" {
			return fmt.Errorf("--version is required")
		}

		if airgapOutput == "" {
			airgapOutput = fmt.Sprintf("security-platform-airgap-%s.tar.gz", airgapVersion)
		}

		fmt.Printf("Generating airgap bundle version %s...\n", airgapVersion)
		fmt.Printf("Output: %s\n", airgapOutput)

		// Create bundle structure
		bundleDir := fmt.Sprintf("security-platform-airgap-%s", airgapVersion)
		if err := os.MkdirAll(bundleDir, 0755); err != nil {
			return fmt.Errorf("failed to create bundle directory: %w", err)
		}
		defer os.RemoveAll(bundleDir)

		// Collect all components
		if err := collectAgents(bundleDir); err != nil {
			return err
		}
		if err := collectSDKs(bundleDir); err != nil {
			return err
		}
		if err := collectInstallers(bundleDir); err != nil {
			return err
		}
		if err := collectDependencies(bundleDir); err != nil {
			return err
		}
		if err := createManifest(bundleDir, airgapVersion); err != nil {
			return err
		}
		if err := createInstallScript(bundleDir); err != nil {
			return err
		}

		// Create tar.gz archive
		if err := createArchive(bundleDir, airgapOutput); err != nil {
			return fmt.Errorf("failed to create archive: %w", err)
		}

		fmt.Printf("Airgap bundle created: %s\n", airgapOutput)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(airgapCmd)
	airgapCmd.Flags().StringVar(&airgapOutput, "output", "", "Output file path (default: security-platform-airgap-<version>.tar.gz)")
	airgapCmd.Flags().StringVar(&airgapVersion, "version", "", "Version number (required)")
	airgapCmd.Flags().StringSliceVar(&airgapIncludeAgents, "agents", []string{}, "Specific agents to include (default: all)")
	airgapCmd.Flags().StringSliceVar(&airgapIncludePlatforms, "platforms", []string{}, "Specific platforms to include (default: all)")
}

func collectAgents(bundleDir string) error {
	agentsDir := filepath.Join(bundleDir, "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		return err
	}

	// Collect all agent binaries
	agents := []string{"java", "dotnet", "go", "node", "python"}
	for _, agent := range agents {
		agentPath := filepath.Join("agents", agent)
		if _, err := os.Stat(agentPath); err == nil {
			destPath := filepath.Join(agentsDir, agent)
			if err := copyDir(agentPath, destPath); err != nil {
				return fmt.Errorf("failed to copy agent %s: %w", agent, err)
			}
		}
	}

	return nil
}

func collectSDKs(bundleDir string) error {
	sdksDir := filepath.Join(bundleDir, "sdks")
	if err := os.MkdirAll(sdksDir, 0755); err != nil {
		return err
	}

	// Collect mobile SDKs
	sdks := []string{"ios", "android", "frontend"}
	for _, sdk := range sdks {
		sdkPath := filepath.Join("agents", sdk)
		if _, err := os.Stat(sdkPath); err == nil {
			destPath := filepath.Join(sdksDir, sdk)
			if err := copyDir(sdkPath, destPath); err != nil {
				return fmt.Errorf("failed to copy SDK %s: %w", sdk, err)
			}
		}
	}

	return nil
}

func collectInstallers(bundleDir string) error {
	installersDir := filepath.Join(bundleDir, "installers")
	if err := os.MkdirAll(installersDir, 0755); err != nil {
		return err
	}

	// Collect package installers
	installers := []string{"homebrew", "chocolatey", "deb", "rpm", "msi", "dmg"}
	for _, installer := range installers {
		installerPath := filepath.Join("packaging", installer)
		if _, err := os.Stat(installerPath); err == nil {
			destPath := filepath.Join(installersDir, installer)
			if err := copyDir(installerPath, destPath); err != nil {
				return fmt.Errorf("failed to copy installer %s: %w", installer, err)
			}
		}
	}

	return nil
}

func collectDependencies(bundleDir string) error {
	depsDir := filepath.Join(bundleDir, "dependencies")
	if err := os.MkdirAll(depsDir, 0755); err != nil {
		return err
	}

	// Create dependency manifest
	manifest := `# Dependencies Manifest
# This file lists all external dependencies required for offline installation

## Java Agent
- OpenTelemetry Java SDK
- Jackson YAML

## .NET Agent
- OpenTelemetry .NET SDK
- Microsoft.Extensions.*

## Go Agent
- OpenTelemetry Go SDK
- zap logger

## Node.js Agent
- @opentelemetry/sdk-trace-node
- @opentelemetry/exporter-otlp-http

## Python Agent
- opentelemetry-api
- opentelemetry-sdk

## iOS SDK
- OpenTelemetry Swift SDK

## Android SDK
- OpenTelemetry Android SDK
- Timber logging

## Kubernetes Operator
- Kubernetes client-go
- controller-runtime
`

	if err := os.WriteFile(filepath.Join(depsDir, "MANIFEST.md"), []byte(manifest), 0644); err != nil {
		return err
	}

	return nil
}

func createManifest(bundleDir, version string) error {
	manifest := fmt.Sprintf(`# Security Platform Airgap Bundle

Version: %s
Generated: %s

## Contents

- agents/ - All agent binaries and packages
- sdks/ - Mobile and frontend SDKs
- installers/ - Package installers (deb, rpm, msi, dmg, homebrew, chocolatey)
- dependencies/ - Dependency manifests
- install.sh - Installation script

## Installation

Extract the bundle and run:

    ./install.sh

Or follow platform-specific installation instructions in the installers/ directory.
`, version, "now")

	return os.WriteFile(filepath.Join(bundleDir, "README.md"), []byte(manifest), 0644)
}

func createInstallScript(bundleDir string) error {
	script := `#!/bin/bash
# Security Platform Airgap Installation Script

set -e

echo "Security Platform Airgap Installation"
echo "======================================"

# Detect platform
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    PLATFORM="linux"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    PLATFORM="macos"
elif [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
    PLATFORM="windows"
else
    echo "Unsupported platform: $OSTYPE"
    exit 1
fi

echo "Detected platform: $PLATFORM"

# Install based on platform
case $PLATFORM in
    linux)
        if command -v apt-get &> /dev/null; then
            echo "Installing DEB package..."
            sudo dpkg -i installers/deb/*.deb
        elif command -v yum &> /dev/null || command -v dnf &> /dev/null; then
            echo "Installing RPM package..."
            sudo rpm -i installers/rpm/*.rpm
        else
            echo "Please install manually from installers/"
        fi
        ;;
    macos)
        if command -v brew &> /dev/null; then
            echo "Installing via Homebrew..."
            brew install installers/homebrew/*.rb
        else
            echo "Please install manually from installers/"
        fi
        ;;
    windows)
        echo "Please install manually using installers/chocolatey/"
        ;;
esac

echo "Installation complete!"
`

	scriptPath := filepath.Join(bundleDir, "install.sh")
	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		return err
	}
	return nil
}

func createArchive(sourceDir, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	gzw := gzip.NewWriter(file)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath := strings.TrimPrefix(path, sourceDir)
		relPath = strings.TrimPrefix(relPath, string(filepath.Separator))

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = relPath

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
		}

		return nil
	})
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath := strings.TrimPrefix(path, src)
		relPath = strings.TrimPrefix(relPath, string(filepath.Separator))
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}
