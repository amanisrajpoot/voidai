package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var buildLang string
var buildVersion string
var buildArch string

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build agent",
	Short: "Compile and produce agent artifact(s)",
	Long: `Build an agent for the specified language and version.
Produces compiled artifacts ready for packaging.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if args[0] != "agent" {
			return fmt.Errorf("only 'agent' build target is supported")
		}

		if buildLang == "" {
			return fmt.Errorf("--lang is required")
		}

		if buildVersion == "" {
			return fmt.Errorf("--version is required")
		}

		fmt.Printf("Building %s agent version %s\n", buildLang, buildVersion)

		outDir := viper.GetString("out")
		if err := os.MkdirAll(outDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		// Build based on language
		switch buildLang {
		case "node":
			return buildNodeAgent(outDir, buildVersion)
		case "python":
			return buildPythonAgent(outDir, buildVersion)
		case "java":
			return buildJavaAgent(outDir, buildVersion)
		case "dotnet":
			return buildDotNetAgent(outDir, buildVersion)
		case "go":
			return buildGoAgent(outDir, buildVersion, buildArch)
		case "frontend":
			return buildFrontendSDK(outDir, buildVersion)
		default:
			return fmt.Errorf("unsupported language: %s", buildLang)
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVar(&buildLang, "lang", "", "language/runtime (node, python, java, dotnet, go, frontend)")
	buildCmd.Flags().StringVar(&buildVersion, "version", "", "version number (e.g., 1.2.0)")
	buildCmd.Flags().StringVar(&buildArch, "arch", "amd64", "architecture (amd64, arm64, 386)")
}

func buildNodeAgent(outDir, version string) error {
	fmt.Println("Building Node.js agent...")
	agentDir := filepath.Join("..", "agents", "node")
	
	// Check if agent directory exists
	if _, err := os.Stat(agentDir); os.IsNotExist(err) {
		return fmt.Errorf("agent directory not found: %s", agentDir)
	}
	
	// Create build instructions file
	buildScript := fmt.Sprintf(`#!/bin/bash
set -e
cd %s
npm install
npm run build
cp -r dist %s/node-agent-dist
`, agentDir, outDir)
	
	scriptPath := filepath.Join(outDir, "build-node.sh")
	if err := os.WriteFile(scriptPath, []byte(buildScript), 0755); err != nil {
		return err
	}
	
	fmt.Printf("  Created build script: %s\n", scriptPath)
	fmt.Println("  Run: bash build-node.sh to build")
	return nil
}

func buildPythonAgent(outDir, version string) error {
	fmt.Println("Building Python agent...")
	agentDir := filepath.Join("..", "agents", "python")
	
	if _, err := os.Stat(agentDir); os.IsNotExist(err) {
		return fmt.Errorf("agent directory not found: %s", agentDir)
	}
	
	buildScript := fmt.Sprintf(`#!/bin/bash
set -e
cd %s
python3 -m pip install --upgrade build
python3 -m build
cp dist/* %s/
`, agentDir, outDir)
	
	scriptPath := filepath.Join(outDir, "build-python.sh")
	if err := os.WriteFile(scriptPath, []byte(buildScript), 0755); err != nil {
		return err
	}
	
	fmt.Printf("  Created build script: %s\n", scriptPath)
	fmt.Println("  Run: bash build-python.sh to build")
	return nil
}

func buildJavaAgent(outDir, version string) error {
	fmt.Println("Building Java agent...")
	agentDir := filepath.Join("..", "agents", "java")
	
	if _, err := os.Stat(agentDir); os.IsNotExist(err) {
		return fmt.Errorf("agent directory not found: %s", agentDir)
	}
	
	buildScript := fmt.Sprintf(`#!/bin/bash
set -e
cd %s
mvn clean package -DskipTests
cp target/*.jar %s/
`, agentDir, outDir)
	
	scriptPath := filepath.Join(outDir, "build-java.sh")
	if err := os.WriteFile(scriptPath, []byte(buildScript), 0755); err != nil {
		return err
	}
	
	fmt.Printf("  Created build script: %s\n", scriptPath)
	fmt.Println("  Run: bash build-java.sh to build")
	return nil
}

func buildDotNetAgent(outDir, version string) error {
	fmt.Println("Building .NET agent...")
	agentDir := filepath.Join("..", "agents", "dotnet", "SecurityPlatform.Agent")
	
	if _, err := os.Stat(agentDir); os.IsNotExist(err) {
		return fmt.Errorf("agent directory not found: %s", agentDir)
	}
	
	buildScript := fmt.Sprintf(`#!/bin/bash
set -e
cd %s
dotnet build --configuration Release
dotnet pack --configuration Release --output %s
`, agentDir, outDir)
	
	scriptPath := filepath.Join(outDir, "build-dotnet.sh")
	if err := os.WriteFile(scriptPath, []byte(buildScript), 0755); err != nil {
		return err
	}
	
	fmt.Printf("  Created build script: %s\n", scriptPath)
	fmt.Println("  Run: bash build-dotnet.sh to build")
	return nil
}

func buildGoAgent(outDir, version, arch string) error {
	fmt.Printf("Building Go agent for %s...\n", arch)
	agentDir := filepath.Join("..", "agents", "go")
	
	if _, err := os.Stat(agentDir); os.IsNotExist(err) {
		return fmt.Errorf("agent directory not found: %s", agentDir)
	}
	
	outputName := fmt.Sprintf("security-platform-agent-%s-%s", runtime.GOOS, arch)
	if runtime.GOOS == "windows" {
		outputName += ".exe"
	}
	
	buildScript := fmt.Sprintf(`#!/bin/bash
set -e
cd %s
GOOS=%s GOARCH=%s go build -o %s/%s .
`, agentDir, runtime.GOOS, arch, outDir, outputName)
	
	scriptPath := filepath.Join(outDir, "build-go.sh")
	if err := os.WriteFile(scriptPath, []byte(buildScript), 0755); err != nil {
		return err
	}
	
	fmt.Printf("  Created build script: %s\n", scriptPath)
	fmt.Printf("  Output: %s/%s\n", outDir, outputName)
	fmt.Println("  Run: bash build-go.sh to build")
	return nil
}

func buildFrontendSDK(outDir, version string) error {
	fmt.Println("Building Frontend SDK...")
	agentDir := filepath.Join("..", "agents", "frontend")
	
	if _, err := os.Stat(agentDir); os.IsNotExist(err) {
		return fmt.Errorf("agent directory not found: %s", agentDir)
	}
	
	buildScript := fmt.Sprintf(`#!/bin/bash
set -e
cd %s
npm install
npm run build
cp -r dist %s/frontend-sdk-dist
`, agentDir, outDir)
	
	scriptPath := filepath.Join(outDir, "build-frontend.sh")
	if err := os.WriteFile(scriptPath, []byte(buildScript), 0755); err != nil {
		return err
	}
	
	fmt.Printf("  Created build script: %s\n", scriptPath)
	fmt.Println("  Run: bash build-frontend.sh to build")
	return nil
}
