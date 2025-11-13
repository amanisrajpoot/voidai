package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/observability-platform/builder/internal/scaffold"
	"github.com/spf13/cobra"
)

var lang string
var projectName string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new agent project",
	Long: `Initialize a new agent project with scaffolding for the specified language.

Supported languages:
  - node (Node.js/Express)
  - python (Python/FastAPI/Django)
  - java (Java/Spring Boot)
  - dotnet (.NET/C#)
  - go (Go)
  - ruby (Ruby/Rails)
  - frontend (TypeScript/React/Vue/Angular)
  - ios (Swift/iOS)
  - android (Kotlin/Android)
  - desktop-windows (C# Windows Service)
  - desktop-macos (Swift macOS Daemon)
  - desktop-linux (Go Linux Daemon)`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName = args[0]
		
		if lang == "" {
			return fmt.Errorf("--lang is required")
		}

		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		projectPath := filepath.Join(wd, projectName)
		
		if err := os.MkdirAll(projectPath, 0755); err != nil {
			return fmt.Errorf("failed to create project directory: %w", err)
		}

		scaffolder := scaffold.NewScaffolder(projectPath, lang)
		if err := scaffolder.Scaffold(); err != nil {
			return fmt.Errorf("failed to scaffold project: %w", err)
		}

		fmt.Printf("✓ Successfully initialized %s agent project in %s\n", lang, projectPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&lang, "lang", "", "Language/runtime (required)")
}
