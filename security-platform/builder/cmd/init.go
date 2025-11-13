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
	Long: `Scaffold a new agent project template for the chosen runtime.
Supports: node, python, java, dotnet, go, ios, android`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		lang, _ := cmd.Flags().GetString("lang")
		
		if lang == "" {
			return fmt.Errorf("--lang is required (node, python, java, dotnet, go, ios, android)")
		}

		projectDir := filepath.Join(".", projectName)
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			return fmt.Errorf("failed to create project directory: %w", err)
		}

		return scaffoldProject(projectDir, lang, projectName)
	},
}

func init() {
	initCmd.Flags().String("lang", "", "Language/runtime (node, python, java, dotnet, go, ios, android)")
	rootCmd.AddCommand(initCmd)
}

func scaffoldProject(dir, lang, name string) error {
	// Create basic project structure
	config := map[string]interface{}{
		"name":    name,
		"version": "0.1.0",
		"lang":    lang,
	}

	// Write config file
	configPath := filepath.Join(dir, "agent.yaml")
	if err := writeConfig(configPath, config); err != nil {
		return err
	}

	// Create language-specific files
	switch lang {
	case "node":
		return scaffoldNode(dir, name)
	case "python":
		return scaffoldPython(dir, name)
	case "java":
		return scaffoldJava(dir, name)
	case "dotnet":
		return scaffoldDotnet(dir, name)
	case "go":
		return scaffoldGo(dir, name)
	case "ios":
		return scaffoldIOS(dir, name)
	case "android":
		return scaffoldAndroid(dir, name)
	default:
		return fmt.Errorf("unsupported language: %s", lang)
	}
}

func scaffoldNode(dir, name string) error {
	packageJson := `{
  "name": "` + name + `",
  "version": "0.1.0",
  "main": "index.js",
  "dependencies": {
    "@opentelemetry/api": "^1.7.0",
    "@opentelemetry/sdk-trace-node": "^1.20.0",
    "@opentelemetry/exporter-otlp-http": "^0.45.0"
  }
}`
	return os.WriteFile(filepath.Join(dir, "package.json"), []byte(packageJson), 0644)
}

func scaffoldPython(dir, name string) error {
	setupPy := `from setuptools import setup, find_packages

setup(
    name="` + name + `",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[
        "opentelemetry-api>=1.20.0",
        "opentelemetry-sdk>=1.20.0",
        "opentelemetry-exporter-otlp-proto-http>=1.20.0",
    ],
)
`
	return os.WriteFile(filepath.Join(dir, "setup.py"), []byte(setupPy), 0644)
}

func scaffoldJava(dir, name string) error {
	// Create Maven structure
	pomPath := filepath.Join(dir, "pom.xml")
	pom := `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 
         http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.securityplatform</groupId>
    <artifactId>` + name + `</artifactId>
    <version>0.1.0</version>
    <dependencies>
        <dependency>
            <groupId>io.opentelemetry</groupId>
            <artifactId>opentelemetry-api</artifactId>
            <version>1.32.0</version>
        </dependency>
    </dependencies>
</project>
`
	return os.WriteFile(pomPath, []byte(pom), 0644)
}

func scaffoldDotnet(dir, name string) error {
	csproj := `<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <TargetFramework>net8.0</TargetFramework>
    <RootNamespace>` + name + `</RootNamespace>
  </PropertyGroup>
  <ItemGroup>
    <PackageReference Include="OpenTelemetry" Version="1.7.0" />
  </ItemGroup>
</Project>
`
	return os.WriteFile(filepath.Join(dir, name+".csproj"), []byte(csproj), 0644)
}

func scaffoldGo(dir, name string) error {
	goMod := `module github.com/security-platform/` + name + `

go 1.21

require (
	go.opentelemetry.io/otel v1.21.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/http v1.21.0
)
`
	return os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644)
}

func scaffoldIOS(dir, name string) error {
	// Create Podspec
	podspec := `Pod::Spec.new do |s|
  s.name         = "` + name + `"
  s.version      = "0.1.0"
  s.summary      = "Security Platform iOS SDK"
  s.homepage     = "https://github.com/security-platform"
  s.license      = "MIT"
  s.author       = { "Security Platform" => "support@securityplatform.com" }
  s.platform     = :ios, "13.0"
  s.source       = { :git => "", :tag => "#{s.version}" }
  s.source_files = "Sources/**/*.swift"
end
`
	return os.WriteFile(filepath.Join(dir, name+".podspec"), []byte(podspec), 0644)
}

func scaffoldAndroid(dir, name string) error {
	// Create build.gradle
	buildGradle := `plugins {
    id 'com.android.library'
    id 'org.jetbrains.kotlin.android'
}

android {
    namespace 'com.securityplatform.` + name + `'
    compileSdk 34
    defaultConfig {
        minSdk 24
    }
}

dependencies {
    implementation 'io.opentelemetry:opentelemetry-api:1.32.0'
}
`
	return os.WriteFile(filepath.Join(dir, "build.gradle.kts"), []byte(buildGradle), 0644)
}

func writeConfig(path string, config map[string]interface{}) error {
	// Simple YAML writer (in production, use proper YAML library)
	yaml := fmt.Sprintf(`name: %s
version: %s
lang: %s
`, config["name"], config["version"], config["lang"])
	return os.WriteFile(path, []byte(yaml), 0644)
}
