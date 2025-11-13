package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Package artifacts into installers",
	Long: `Produce installers for various platforms.
Supported targets: deb, rpm, msi, dmg, helm, ova`,
	RunE: func(cmd *cobra.Command, args []string) error {
		targets, _ := cmd.Flags().GetString("target")
		outDir, _ := cmd.Flags().GetString("out")
		version, _ := cmd.Flags().GetString("version")

		if targets == "" {
			return fmt.Errorf("--target is required (comma-separated: deb,rpm,msi,dmg,helm)")
		}

		if version == "" {
			version = "0.1.0"
		}

		targetList := strings.Split(targets, ",")
		for _, target := range targetList {
			target = strings.TrimSpace(target)
			if err := packageForTarget(target, version, outDir); err != nil {
				return fmt.Errorf("failed to package for %s: %w", target, err)
			}
		}

		return nil
	},
}

func init() {
	packageCmd.Flags().String("target", "", "Target platforms (comma-separated: deb,rpm,msi,dmg,helm)")
	packageCmd.Flags().String("out", "./dist", "Output directory")
	packageCmd.Flags().String("version", "", "Version number")
	rootCmd.AddCommand(packageCmd)
}

func packageForTarget(target, version, outDir string) error {
	fmt.Printf("Packaging for %s...\n", target)

	switch target {
	case "deb":
		return packageDeb(version, outDir)
	case "rpm":
		return packageRpm(version, outDir)
	case "msi":
		return packageMsi(version, outDir)
	case "dmg":
		return packageDmg(version, outDir)
	case "helm":
		return packageHelm(version, outDir)
	case "ova":
		return packageOva(version, outDir)
	default:
		return fmt.Errorf("unsupported target: %s", target)
	}
}

func packageDeb(version, outDir string) error {
	// Create DEB package structure
	debDir := filepath.Join(outDir, "deb")
	controlDir := filepath.Join(debDir, "DEBIAN")
	if err := os.MkdirAll(controlDir, 0755); err != nil {
		return err
	}

	// Write control file
	control := fmt.Sprintf(`Package: security-platform-agent
Version: %s
Section: utils
Priority: optional
Architecture: amd64
Maintainer: Security Platform <support@securityplatform.com>
Description: Security Platform Agent
`, version)
	
	if err := os.WriteFile(filepath.Join(controlDir, "control"), []byte(control), 0644); err != nil {
		return err
	}

	// Create binary location
	binaryDir := filepath.Join(debDir, "usr", "local", "bin")
	if err := os.MkdirAll(binaryDir, 0755); err != nil {
		return err
	}

	// Build DEB package
	cmd := exec.Command("dpkg-deb", "--build", debDir, filepath.Join(outDir, fmt.Sprintf("security-platform-agent_%s_amd64.deb", version)))
	return cmd.Run()
}

func packageRpm(version, outDir string) error {
	// Create RPM spec file
	spec := fmt.Sprintf(`Name: security-platform-agent
Version: %s
Release: 1
Summary: Security Platform Agent
License: MIT
BuildArch: x86_64

%%description
Security Platform Agent

%%install
mkdir -p %%{buildroot}/usr/local/bin
cp agent %%{buildroot}/usr/local/bin/

%%files
/usr/local/bin/agent
`, version)

	specPath := filepath.Join(outDir, "agent.spec")
	if err := os.WriteFile(specPath, []byte(spec), 0644); err != nil {
		return err
	}

	// Build RPM (requires rpmbuild)
	cmd := exec.Command("rpmbuild", "-bb", specPath)
	return cmd.Run()
}

func packageMsi(version, outDir string) error {
	// MSI requires WiX Toolset
	// Create a basic WiX XML file
	wix := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<Wix xmlns="http://schemas.microsoft.com/wix/2006/wi">
  <Product Id="*" Name="Security Platform Agent" Language="1033" Version="%s" Manufacturer="Security Platform" UpgradeCode="GUID-HERE">
    <Package InstallerVersion="200" Compressed="yes" InstallScope="perMachine" />
    <MajorUpgrade DowngradeErrorMessage="A newer version is already installed." />
    <MediaTemplate />
    <Feature Id="ProductFeature" Title="Security Platform Agent" Level="1">
      <ComponentRef Id="AgentBinary" />
    </Feature>
  </Product>
  <Fragment>
    <Directory Id="TARGETDIR" Name="SourceDir">
      <Directory Id="ProgramFiles64Folder">
        <Directory Id="INSTALLFOLDER" Name="SecurityPlatform">
          <Component Id="AgentBinary" Guid="GUID-HERE">
            <File Id="AgentExe" Source="agent.exe" KeyPath="yes" />
          </Component>
        </Directory>
      </Directory>
    </Directory>
  </Fragment>
</Wix>
`, version)

	wixPath := filepath.Join(outDir, "agent.wxs")
	if err := os.WriteFile(wixPath, []byte(wix), 0644); err != nil {
		return err
	}

	fmt.Println("WiX file created. Run 'candle' and 'light' to build MSI.")
	return nil
}

func packageDmg(version, outDir string) error {
	// DMG creation requires macOS and hdiutil
	appDir := filepath.Join(outDir, "SecurityPlatformAgent.app")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return err
	}

	// Create DMG
	dmgPath := filepath.Join(outDir, fmt.Sprintf("security-platform-agent-%s.dmg", version))
	cmd := exec.Command("hdiutil", "create", "-volname", "Security Platform Agent", "-srcfolder", appDir, "-ov", "-format", "UDZO", dmgPath)
	return cmd.Run()
}

func packageHelm(version, outDir string) error {
	helmDir := filepath.Join(outDir, "helm", "security-platform-agent")
	if err := os.MkdirAll(helmDir, 0755); err != nil {
		return err
	}

	// Create Chart.yaml
	chart := fmt.Sprintf(`apiVersion: v2
name: security-platform-agent
description: Security Platform Agent Helm Chart
type: application
version: %s
appVersion: "%s"
`, version, version)

	if err := os.WriteFile(filepath.Join(helmDir, "Chart.yaml"), []byte(chart), 0644); err != nil {
		return err
	}

	// Create values.yaml
	values := `replicaCount: 1

image:
  repository: security-platform/agent
  tag: "latest"
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 8080

config:
  controlPlaneUrl: ""
  authKey: ""
`

	if err := os.WriteFile(filepath.Join(helmDir, "values.yaml"), []byte(values), 0644); err != nil {
		return err
	}

	// Create templates directory
	templatesDir := filepath.Join(helmDir, "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return err
	}

	fmt.Println("Helm chart created at", helmDir)
	return nil
}

func packageOva(version, outDir string) error {
	// OVA creation requires VM tools
	fmt.Println("OVA packaging requires VM export tools. Creating manifest...")
	
	manifest := fmt.Sprintf(`{
  "version": "%s",
  "format": "ova",
  "description": "Security Platform Appliance"
}`, version)

	return os.WriteFile(filepath.Join(outDir, "manifest.json"), []byte(manifest), 0644)
}
