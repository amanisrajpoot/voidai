package cmd

import (
	"fmt"

	"github.com/observability-platform/builder/internal/build"
	"github.com/spf13/cobra"
)

var buildLang string
var buildVersion string
var buildArch string
var buildOS string

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [agent|sdk]",
	Short: "Build an agent or SDK",
	Long: `Build an agent or SDK for the specified language and platform.

Examples:
  builder build agent --lang=node --version=1.2.0
  builder build agent --lang=python --version=1.2.0 --arch=amd64 --os=linux
  builder build sdk --lang=frontend --version=1.2.0`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		
		if target != "agent" && target != "sdk" {
			return fmt.Errorf("target must be 'agent' or 'sdk'")
		}

		if buildLang == "" {
			return fmt.Errorf("--lang is required")
		}

		if buildVersion == "" {
			return fmt.Errorf("--version is required")
		}

		builder := build.NewBuilder(buildLang, buildVersion)
		builder.SetArch(buildArch)
		builder.SetOS(buildOS)

		artifact, err := builder.Build(target)
		if err != nil {
			return fmt.Errorf("build failed: %w", err)
		}

		fmt.Printf("✓ Successfully built %s for %s\n", target, buildLang)
		fmt.Printf("  Artifact: %s\n", artifact)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVar(&buildLang, "lang", "", "Language/runtime (required)")
	buildCmd.Flags().StringVar(&buildVersion, "version", "", "Version (required)")
	buildCmd.Flags().StringVar(&buildArch, "arch", "", "Target architecture (amd64, arm64, etc.)")
	buildCmd.Flags().StringVar(&buildOS, "os", "", "Target OS (linux, darwin, windows, etc.)")
}
