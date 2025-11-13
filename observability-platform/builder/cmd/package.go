package cmd

import (
	"fmt"
	"strings"

	"github.com/observability-platform/builder/internal/package"
	"github.com/spf13/cobra"
)

var packageTargets string
var packageOut string
var packageVersion string

// packageCmd represents the package command
var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Package artifacts into installers",
	Long: `Package built artifacts into platform-specific installers.

Supported targets:
  - deb (Debian/Ubuntu .deb packages)
  - rpm (RedHat/CentOS .rpm packages)
  - msi (Windows MSI installers)
  - dmg (macOS disk images)
  - helm (Kubernetes Helm charts)
  - ova (VMware OVA images)
  - docker (Docker images)
  - homebrew (Homebrew tap)
  - chocolatey (Chocolatey packages)

Examples:
  builder package --target=deb,rpm --out=./dist --version=1.2.0
  builder package --target=msi,dmg --out=./dist --version=1.2.0`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if packageTargets == "" {
			return fmt.Errorf("--target is required")
		}

		if packageOut == "" {
			packageOut = "./dist"
		}

		if packageVersion == "" {
			return fmt.Errorf("--version is required")
		}

		targets := strings.Split(packageTargets, ",")
		for i, t := range targets {
			targets[i] = strings.TrimSpace(t)
		}

		packager := package.NewPackager(packageOut, packageVersion)
		
		artifacts, err := packager.Package(targets)
		if err != nil {
			return fmt.Errorf("packaging failed: %w", err)
		}

		fmt.Printf("✓ Successfully packaged %d installer(s)\n", len(artifacts))
		for _, artifact := range artifacts {
			fmt.Printf("  - %s\n", artifact)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(packageCmd)
	packageCmd.Flags().StringVar(&packageTargets, "target", "", "Comma-separated list of package targets (required)")
	packageCmd.Flags().StringVar(&packageOut, "out", "./dist", "Output directory")
	packageCmd.Flags().StringVar(&packageVersion, "version", "", "Version (required)")
}
