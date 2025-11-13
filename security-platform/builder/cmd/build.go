package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

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
	// In real implementation, run: npm install && npm run build
	// Copy dist/ to outDir
	return nil
}

func buildPythonAgent(outDir, version string) error {
	fmt.Println("Building Python agent...")
	// In real implementation, run: python setup.py sdist bdist_wheel
	return nil
}

func buildJavaAgent(outDir, version string) error {
	fmt.Println("Building Java agent...")
	agentDir := "agents/java"
	if _, err := os.Stat(agentDir); os.IsNotExist(err) {
		return fmt.Errorf("Java agent directory not found: %s", agentDir)
	}
	
	// Run Maven build
	cmd := exec.Command("mvn", "clean", "package", "-DskipTests")
	cmd.Dir = agentDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Maven build failed: %w", err)
	}
	
	// Copy JAR to output
	targetJar := filepath.Join(agentDir, "target", fmt.Sprintf("security-platform-agent-%s.jar", version))
	destJar := filepath.Join(outDir, fmt.Sprintf("security-platform-agent-%s.jar", version))
	if err := copyFile(targetJar, destJar); err != nil {
		return fmt.Errorf("Failed to copy JAR: %w", err)
	}
	
	return nil
}

func buildDotNetAgent(outDir, version string) error {
	fmt.Println("Building .NET agent...")
	agentDir := "agents/dotnet/SecurityPlatform.Agent"
	if _, err := os.Stat(agentDir); os.IsNotExist(err) {
		return fmt.Errorf(".NET agent directory not found: %s", agentDir)
	}
	
	// Run dotnet build
	cmd := exec.Command("dotnet", "build", "--configuration", "Release", "-p:Version="+version)
	cmd.Dir = agentDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("dotnet build failed: %w", err)
	}
	
	// Run dotnet pack
	cmd = exec.Command("dotnet", "pack", "--configuration", "Release", "-p:Version="+version, "--output", outDir)
	cmd.Dir = agentDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("dotnet pack failed: %w", err)
	}
	
	return nil
}

func buildGoAgent(outDir, version, arch string) error {
	fmt.Printf("Building Go agent for %s...\n", arch)
	agentDir := "agents/go"
	if _, err := os.Stat(agentDir); os.IsNotExist(err) {
		return fmt.Errorf("Go agent directory not found: %s", agentDir)
	}
	
	// Set GOOS and GOARCH based on arch
	goos := "linux"
	goarch := "amd64"
	if arch == "arm64" {
		goarch = "arm64"
	}
	
	// Build binary
	binaryName := fmt.Sprintf("security-platform-agent-%s-%s-%s", version, goos, goarch)
	if goos == "windows" {
		binaryName += ".exe"
	}
	
	cmd := exec.Command("go", "build", "-o", filepath.Join(outDir, binaryName), "-ldflags", fmt.Sprintf("-X main.version=%s", version), "./main.go")
	cmd.Dir = agentDir
	cmd.Env = append(os.Environ(), "GOOS="+goos, "GOARCH="+goarch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build failed: %w", err)
	}
	
	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()
	
	_, err = io.Copy(destFile, sourceFile)
	return err
}

func buildFrontendSDK(outDir, version string) error {
	fmt.Println("Building Frontend SDK...")
	// In real implementation, run: npm install && npm run build
	// Create UMD bundle and minified versions
	return nil
}
