package cmd

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var airgapOut string
var airgapVersion string
var airgapInclude string

// airgapCmd represents the airgap command
var airgapCmd = &cobra.Command{
	Use:   "airgap",
	Short: "Generate air-gapped installation bundle",
	Long: `Generate a complete air-gapped installation bundle containing
all necessary packages, dependencies, and installation scripts for offline deployment.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if airgapOut == "" {
			airgapOut = "./airgap-bundle"
		}

		if airgapVersion == "" {
			airgapVersion = "1.0.0"
		}

		fmt.Printf("Generating airgap bundle version %s to %s\n", airgapVersion, airgapOut)

		// Create bundle directory
		if err := os.MkdirAll(airgapOut, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		// Create manifest
		manifest := createManifest(airgapVersion, strings.Split(airgapInclude, ","))
		manifestPath := filepath.Join(airgapOut, "manifest.json")
		manifestData, _ := json.MarshalIndent(manifest, "", "  ")
		if err := os.WriteFile(manifestPath, manifestData, 0644); err != nil {
			return fmt.Errorf("failed to write manifest: %w", err)
		}

		// Create installation script
		if err := createInstallScript(airgapOut); err != nil {
			return fmt.Errorf("failed to create install script: %w", err)
		}

		// Create README
		if err := createReadme(airgapOut, airgapVersion); err != nil {
			return fmt.Errorf("failed to create README: %w", err)
		}

		// Package everything into tar.gz
		archivePath := fmt.Sprintf("%s/security-platform-agent-airgap-%s.tar.gz", airgapOut, airgapVersion)
		if err := createArchive(airgapOut, archivePath); err != nil {
			return fmt.Errorf("failed to create archive: %w", err)
		}

		fmt.Printf("Airgap bundle created: %s\n", archivePath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(airgapCmd)
	airgapCmd.Flags().StringVar(&airgapOut, "out", "./airgap-bundle", "output directory")
	airgapCmd.Flags().StringVar(&airgapVersion, "version", "1.0.0", "bundle version")
	airgapCmd.Flags().StringVar(&airgapInclude, "include", "deb,rpm,msi,dmg", "comma-separated list of packages to include")
}

type Manifest struct {
	Version     string   `json:"version"`
	GeneratedAt string   `json:"generated_at"`
	Packages    []string  `json:"packages"`
	Dependencies []string `json:"dependencies"`
	Checksums   map[string]string `json:"checksums"`
}

func createManifest(version string, includes []string) Manifest {
	return Manifest{
		Version:     version,
		GeneratedAt: fmt.Sprintf("%d", os.Getpid()), // Placeholder
		Packages:    includes,
		Dependencies: []string{},
		Checksums:   make(map[string]string),
	}
}

func createInstallScript(outDir string) error {
	script := `#!/bin/bash
set -e

echo "Installing Security Platform Agent in air-gapped environment..."

BUNDLE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$BUNDLE_DIR"

# Detect OS
if [ -f /etc/debian_version ]; then
    echo "Detected Debian/Ubuntu system"
    if [ -f security-platform-agent*.deb ]; then
        sudo dpkg -i security-platform-agent*.deb
    else
        echo "ERROR: No .deb package found"
        exit 1
    fi
elif [ -f /etc/redhat-release ]; then
    echo "Detected RHEL/CentOS system"
    if [ -f security-platform-agent*.rpm ]; then
        sudo rpm -ivh security-platform-agent*.rpm
    else
        echo "ERROR: No .rpm package found"
        exit 1
    fi
elif [[ "$OSTYPE" == "darwin"* ]]; then
    echo "Detected macOS system"
    if [ -f security-platform-agent*.pkg ]; then
        sudo installer -pkg security-platform-agent*.pkg -target /
    else
        echo "ERROR: No .pkg package found"
        exit 1
    fi
else
    echo "Unsupported OS. Please install manually."
    exit 1
fi

# Create config directory
sudo mkdir -p /etc/security-platform

# Copy default config if provided
if [ -f config.yaml ]; then
    sudo cp config.yaml /etc/security-platform/config.yaml
fi

# Start service
if command -v systemctl &> /dev/null; then
    sudo systemctl enable security-platform-agent
    sudo systemctl start security-platform-agent
elif command -v launchctl &> /dev/null; then
    sudo launchctl load /Library/LaunchDaemons/com.securityplatform.agent.plist
fi

echo "Installation complete!"
echo ""
echo "Configuration file: /etc/security-platform/config.yaml"
echo "Service status: systemctl status security-platform-agent"
`

	scriptPath := filepath.Join(outDir, "install.sh")
	return os.WriteFile(scriptPath, []byte(script), 0755)
}

func createReadme(outDir string, version string) error {
	readme := fmt.Sprintf(`# Security Platform Agent - Airgap Bundle

Version: %s

## Contents

This bundle contains all necessary files for installing Security Platform Agent
in an air-gapped (offline) environment.

## Installation

1. Extract this bundle:
   \`\`\`bash
   tar -xzf security-platform-agent-airgap-%s.tar.gz
   cd airgap-bundle
   \`\`\`

2. Run the installation script:
   \`\`\`bash
   chmod +x install.sh
   sudo ./install.sh
   \`\`\`

## Manual Installation

If automatic installation fails, you can install packages manually:

### Debian/Ubuntu
\`\`\`bash
sudo dpkg -i security-platform-agent*.deb
\`\`\`

### RHEL/CentOS
\`\`\`bash
sudo rpm -ivh security-platform-agent*.rpm
\`\`\`

### macOS
\`\`\`bash
sudo installer -pkg security-platform-agent*.pkg -target /
\`\`\`

## Configuration

After installation, edit the configuration file:
\`\`\`bash
sudo nano /etc/security-platform/config.yaml
\`\`\`

Key settings:
- \`control_plane_url\`: Your Security Platform control plane URL
- \`auth_key\`: Your authentication key
- \`service_name\`: Name for this agent instance
- \`environment\`: Deployment environment (production, staging, etc.)

## Verification

Check service status:
\`\`\`bash
systemctl status security-platform-agent
\`\`\`

View logs:
\`\`\`bash
journalctl -u security-platform-agent -f
\`\`\`

## Support

For support, visit: https://github.com/security-platform/desktop-agent
`, version, version)

	readmePath := filepath.Join(outDir, "README.md")
	return os.WriteFile(readmePath, []byte(readme), 0644)
}

func createArchive(sourceDir, archivePath string) error {
	file, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzWriter := gzip.NewWriter(file)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the archive itself
		if path == archivePath {
			return nil
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = relPath

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tarWriter, file)
			return err
		}

		return nil
	})
}
