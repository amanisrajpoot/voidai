package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new agent project",
	Long: `Initialize a new agent project with a template for the chosen runtime.
This scaffolds the necessary files and directory structure.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		lang, _ := cmd.Flags().GetString("lang")
		
		if lang == "" {
			return fmt.Errorf("--lang is required (options: python, node, go, java, dotnet, ios, android)")
		}

		workspace, _ := cmd.Flags().GetString("workspace")
		projectPath := filepath.Join(workspace, projectName)

		if err := os.MkdirAll(projectPath, 0755); err != nil {
			return fmt.Errorf("failed to create project directory: %w", err)
		}

		return scaffoldProject(projectPath, lang, projectName)
	},
}

func init() {
	initCmd.Flags().String("lang", "", "Language/runtime (python, node, go, java, dotnet, ios, android)")
	rootCmd.AddCommand(initCmd)
}

func scaffoldProject(path, lang, name string) error {
	// Create basic structure
	dirs := []string{"src", "config", "tests"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(path, dir), 0755); err != nil {
			return err
		}
	}

	// Create agent.yaml config template
	config := fmt.Sprintf(`name: %s
version: 0.1.0
language: %s
control_plane_url: https://api.example.com
auth_key: ""
env: development
redaction_rules:
  - password
  - card
  - ssn
  - pin
  - auth.*
local_policy:
  mode: observe  # observe or block
telemetry_batch_size: 100
otlp_endpoint: http://localhost:4318
`, name, lang)

	if err := os.WriteFile(filepath.Join(path, "agent.yaml"), []byte(config), 0644); err != nil {
		return err
	}

	// Create README
	readme := fmt.Sprintf(`# %s

Agent project for %s runtime.

## Configuration

Edit \`agent.yaml\` to configure:
- Control plane endpoint
- Authentication
- Redaction rules
- Local policy (observe/block)

## Building

\`\`\`bash
builder build agent --lang=%s --version=0.1.0
\`\`\`

## Packaging

\`\`\`bash
builder package --target=deb,rpm --out=./dist
\`\`\`
`, name, lang, lang)

	if err := os.WriteFile(filepath.Join(path, "README.md"), []byte(readme), 0644); err != nil {
		return err
	}

	fmt.Printf("✓ Initialized %s project at %s\n", lang, path)
	return nil
}
