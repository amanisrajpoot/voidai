package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build agent",
	Short: "Build an agent for the specified language",
	Long: `Compile and produce artifacts for the specified agent language.
Supports: node, python, java, dotnet, go`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if args[0] != "agent" {
			return fmt.Errorf("only 'agent' build target is supported")
		}

		lang, _ := cmd.Flags().GetString("lang")
		version, _ := cmd.Flags().GetString("version")
		outDir, _ := cmd.Flags().GetString("out")

		if lang == "" {
			return fmt.Errorf("--lang is required")
		}
		if version == "" {
			version = "0.1.0"
		}

		return buildAgent(lang, version, outDir)
	},
}

func init() {
	buildCmd.Flags().String("lang", "", "Language (node, python, java, dotnet, go)")
	buildCmd.Flags().String("version", "", "Version number")
	buildCmd.Flags().String("out", "./dist", "Output directory")
	rootCmd.AddCommand(buildCmd)
}

func buildAgent(lang, version, outDir string) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	switch lang {
	case "node":
		return buildNodeAgent(version, outDir)
	case "python":
		return buildPythonAgent(version, outDir)
	case "java":
		return buildJavaAgent(version, outDir)
	case "dotnet":
		return buildDotnetAgent(version, outDir)
	case "go":
		return buildGoAgent(version, outDir)
	default:
		return fmt.Errorf("unsupported language: %s", lang)
	}
}

func buildNodeAgent(version, outDir string) error {
	fmt.Println("Building Node.js agent...")
	
	// Run npm pack or build
	cmd := exec.Command("npm", "pack")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("npm pack failed: %w", err)
	}

	// Move tarball to output directory
	glob := fmt.Sprintf("*.tgz")
	matches, _ := filepath.Glob(glob)
	if len(matches) > 0 {
		dest := filepath.Join(outDir, filepath.Base(matches[0]))
		return os.Rename(matches[0], dest)
	}

	return nil
}

func buildPythonAgent(version, outDir string) error {
	fmt.Println("Building Python agent...")
	
	cmd := exec.Command("python", "setup.py", "sdist", "bdist_wheel")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("python build failed: %w", err)
	}

	// Move dist/* to output directory
	distFiles, _ := filepath.Glob("dist/*")
	for _, f := range distFiles {
		dest := filepath.Join(outDir, filepath.Base(f))
		if err := os.Rename(f, dest); err != nil {
			return err
		}
	}

	return nil
}

func buildJavaAgent(version, outDir string) error {
	fmt.Println("Building Java agent...")
	
	cmd := exec.Command("mvn", "clean", "package", "-DskipTests")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("maven build failed: %w", err)
	}

	// Copy JAR to output directory
	jars, _ := filepath.Glob("target/*.jar")
	for _, jar := range jars {
		dest := filepath.Join(outDir, filepath.Base(jar))
		if err := copyFile(jar, dest); err != nil {
			return err
		}
	}

	return nil
}

func buildDotnetAgent(version, outDir string) error {
	fmt.Println("Building .NET agent...")
	
	cmd := exec.Command("dotnet", "pack", "-c", "Release", "-o", outDir)
	cmd.Dir = "."
	return cmd.Run()
}

func buildGoAgent(version, outDir string) error {
	fmt.Println("Building Go agent...")
	
	// Use goreleaser if available, otherwise simple go build
	cmd := exec.Command("go", "build", "-o", filepath.Join(outDir, "agent"), ".")
	cmd.Dir = "."
	return cmd.Run()
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}
