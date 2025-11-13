package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Generate installers and packages",
	Long: `Generate installers and packages for various platforms.
Supports: deb, rpm, msi, dmg, helm, ova, docker`,
	RunE: func(cmd *cobra.Command, args []string) error {
		targets, _ := cmd.Flags().GetString("target")
		if targets == "" {
			return fmt.Errorf("--target is required (e.g., deb,rpm,msi)")
		}

		outDir, _ := cmd.Flags().GetString("out")
		if outDir == "" {
			outDir = "./dist"
		}

		version, _ := cmd.Flags().GetString("version")
		if version == "" {
			return fmt.Errorf("--version is required")
		}

		targetList := strings.Split(targets, ",")
		return generatePackages(targetList, version, outDir)
	},
}

func init() {
	packageCmd.Flags().String("target", "", "Package targets (comma-separated: deb,rpm,msi,dmg,helm,ova)")
	packageCmd.Flags().String("out", "./dist", "Output directory")
	packageCmd.Flags().String("version", "", "Version number")
	rootCmd.AddCommand(packageCmd)
}

func generatePackages(targets []string, version, outDir string) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for _, target := range targets {
		target = strings.TrimSpace(target)
		fmt.Printf("Generating %s package...\n", target)

		var err error
		switch target {
		case "deb":
			err = generateDebPackage(version, outDir)
		case "rpm":
			err = generateRpmPackage(version, outDir)
		case "msi":
			err = generateMsiPackage(version, outDir)
		case "dmg":
			err = generateDmgPackage(version, outDir)
		case "helm":
			err = generateHelmChart(version, outDir)
		case "ova":
			err = generateOvaPackage(version, outDir)
		case "docker":
			err = generateDockerImage(version, outDir)
		default:
			return fmt.Errorf("unsupported package target: %s", target)
		}

		if err != nil {
			return fmt.Errorf("failed to generate %s: %w", target, err)
		}
	}

	return nil
}

func generateDebPackage(version, outDir string) error {
	fmt.Println("  → Creating .deb package...")
	// TODO: Use dpkg-deb or fpm to create .deb
	artifact := filepath.Join(outDir, fmt.Sprintf("security-agent_%s_amd64.deb", version))
	if err := os.WriteFile(artifact, []byte("# DEB package placeholder\n"), 0644); err != nil {
		return err
	}
	fmt.Printf("  ✓ Created: %s\n", artifact)
	return nil
}

func generateRpmPackage(version, outDir string) error {
	fmt.Println("  → Creating .rpm package...")
	// TODO: Use rpmbuild or fpm to create .rpm
	artifact := filepath.Join(outDir, fmt.Sprintf("security-agent-%s-1.x86_64.rpm", version))
	if err := os.WriteFile(artifact, []byte("# RPM package placeholder\n"), 0644); err != nil {
		return err
	}
	fmt.Printf("  ✓ Created: %s\n", artifact)
	return nil
}

func generateMsiPackage(version, outDir string) error {
	fmt.Println("  → Creating .msi package...")
	// TODO: Use WiX Toolset to create .msi
	artifact := filepath.Join(outDir, fmt.Sprintf("security-agent-%s.msi", version))
	if err := os.WriteFile(artifact, []byte("# MSI package placeholder\n"), 0644); err != nil {
		return err
	}
	fmt.Printf("  ✓ Created: %s\n", artifact)
	return nil
}

func generateDmgPackage(version, outDir string) error {
	fmt.Println("  → Creating .dmg package...")
	// TODO: Use hdiutil or create-dmg to create .dmg
	artifact := filepath.Join(outDir, fmt.Sprintf("security-agent-%s.dmg", version))
	if err := os.WriteFile(artifact, []byte("# DMG package placeholder\n"), 0644); err != nil {
		return err
	}
	fmt.Printf("  ✓ Created: %s\n", artifact)
	return nil
}

func generateHelmChart(version, outDir string) error {
	fmt.Println("  → Creating Helm chart...")
	chartDir := filepath.Join(outDir, "security-agent-chart")
	if err := os.MkdirAll(chartDir, 0755); err != nil {
		return err
	}

	// Create Chart.yaml
	chartYaml := fmt.Sprintf(`apiVersion: v2
name: security-agent
description: Security observability platform agent
type: application
version: %s
appVersion: "%s"
`, version, version)
	if err := os.WriteFile(filepath.Join(chartDir, "Chart.yaml"), []byte(chartYaml), 0644); err != nil {
		return err
	}

	// Create values.yaml
	valuesYaml := `controlPlane:
  url: https://api.example.com
  authKey: ""

agent:
  image:
    repository: security-agent
    tag: latest
  resources:
    requests:
      cpu: 100m
      memory: 128Mi
`
	if err := os.WriteFile(filepath.Join(chartDir, "values.yaml"), []byte(valuesYaml), 0644); err != nil {
		return err
	}

	fmt.Printf("  ✓ Created: %s\n", chartDir)
	return nil
}

func generateOvaPackage(version, outDir string) error {
	fmt.Println("  → Creating OVA package...")
	// TODO: Create OVA from VM image
	artifact := filepath.Join(outDir, fmt.Sprintf("security-agent-%s.ova", version))
	if err := os.WriteFile(artifact, []byte("# OVA package placeholder\n"), 0644); err != nil {
		return err
	}
	fmt.Printf("  ✓ Created: %s\n", artifact)
	return nil
}

func generateDockerImage(version, outDir string) error {
	fmt.Println("  → Creating Docker image...")
	// TODO: Build Docker image
	artifact := filepath.Join(outDir, fmt.Sprintf("security-agent-%s.tar", version))
	if err := os.WriteFile(artifact, []byte("# Docker image placeholder\n"), 0644); err != nil {
		return err
	}
	fmt.Printf("  ✓ Created: %s\n", artifact)
	return nil
}
