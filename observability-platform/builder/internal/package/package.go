package package

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Packager struct {
	outDir  string
	version string
}

func NewPackager(outDir, version string) *Packager {
	return &Packager{
		outDir:  outDir,
		version: version,
	}
}

func (p *Packager) Package(targets []string) ([]string, error) {
	if err := os.MkdirAll(p.outDir, 0755); err != nil {
		return nil, err
	}

	var artifacts []string

	for _, target := range targets {
		artifact, err := p.packageTarget(target)
		if err != nil {
			return nil, fmt.Errorf("failed to package %s: %w", target, err)
		}
		artifacts = append(artifacts, artifact)
	}

	return artifacts, nil
}

func (p *Packager) packageTarget(target string) (string, error) {
	switch target {
	case "deb":
		return p.packageDeb()
	case "rpm":
		return p.packageRPM()
	case "msi":
		return p.packageMSI()
	case "dmg":
		return p.packageDMG()
	case "helm":
		return p.packageHelm()
	case "ova":
		return p.packageOVA()
	case "docker":
		return p.packageDocker()
	case "homebrew":
		return p.packageHomebrew()
	case "chocolatey":
		return p.packageChocolatey()
	default:
		return "", fmt.Errorf("unsupported package target: %s", target)
	}
}

func (p *Packager) packageDeb() (string, error) {
	// Create DEB package structure
	debDir := filepath.Join(p.outDir, "deb")
	if err := os.MkdirAll(filepath.Join(debDir, "DEBIAN"), 0755); err != nil {
		return "", err
	}

	// Write control file
	control := fmt.Sprintf(`Package: observability-agent
Version: %s
Architecture: amd64
Maintainer: Observability Platform <support@observability.platform>
Description: Security Observability Agent
`, p.version)

	if err := os.WriteFile(filepath.Join(debDir, "DEBIAN", "control"), []byte(control), 0644); err != nil {
		return "", err
	}

	// Build DEB package
	debFile := filepath.Join(p.outDir, fmt.Sprintf("observability-agent_%s_amd64.deb", p.version))
	cmd := exec.Command("dpkg-deb", "--build", debDir, debFile)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("dpkg-deb failed: %w", err)
	}

	return debFile, nil
}

func (p *Packager) packageRPM() (string, error) {
	// Create RPM spec file
	specFile := filepath.Join(p.outDir, "observability-agent.spec")
	spec := fmt.Sprintf(`Name: observability-agent
Version: %s
Release: 1
Summary: Security Observability Agent
License: Apache-2.0
BuildArch: x86_64

%%description
Security Observability Agent

%%files
`, p.version)

	if err := os.WriteFile(specFile, []byte(spec), 0644); err != nil {
		return "", err
	}

	// Build RPM (requires rpmbuild)
	rpmFile := filepath.Join(p.outDir, fmt.Sprintf("observability-agent-%s-1.x86_64.rpm", p.version))
	cmd := exec.Command("rpmbuild", "-bb", specFile, "--define", fmt.Sprintf("_topdir %s", p.outDir))
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("rpmbuild failed: %w", err)
	}

	return rpmFile, nil
}

func (p *Packager) packageMSI() (string, error) {
	// Create WiX source file
	wxsFile := filepath.Join(p.outDir, "observability-agent.wxs")
	wxs := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<Wix xmlns="http://schemas.microsoft.com/wix/2006/wi">
  <Product Id="*" Name="Observability Agent" Language="1033" Version="%s" Manufacturer="Observability Platform" UpgradeCode="YOUR-GUID-HERE">
    <Package InstallerVersion="200" Compressed="yes" InstallScope="perMachine" />
    <MajorUpgrade DowngradeErrorMessage="A newer version is already installed." />
    <MediaTemplate />
    <Feature Id="ProductFeature" Title="Observability Agent" Level="1">
      <ComponentRef Id="ApplicationFiles" />
    </Feature>
  </Product>
  <Fragment>
    <Directory Id="TARGETDIR" Name="SourceDir">
      <Directory Id="ProgramFilesFolder">
        <Directory Id="INSTALLFOLDER" Name="ObservabilityAgent" />
      </Directory>
    </Directory>
  </Fragment>
  <Fragment>
    <ComponentGroup Id="ProductComponents" Directory="INSTALLFOLDER">
      <Component Id="ApplicationFiles">
        <File Id="AgentExe" Source="observability-agent.exe" />
      </Component>
    </ComponentGroup>
  </Fragment>
</Wix>
`, p.version)

	if err := os.WriteFile(wxsFile, []byte(wxs), 0644); err != nil {
		return "", err
	}

	// Build MSI (requires WiX Toolset)
	msiFile := filepath.Join(p.outDir, fmt.Sprintf("observability-agent-%s.msi", p.version))
	cmd := exec.Command("candle", wxsFile)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("candle failed: %w", err)
	}

	cmd = exec.Command("light", strings.TrimSuffix(wxsFile, ".wxs")+".wixobj", "-out", msiFile)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("light failed: %w", err)
	}

	return msiFile, nil
}

func (p *Packager) packageDMG() (string, error) {
	// Create DMG (macOS)
	dmgFile := filepath.Join(p.outDir, fmt.Sprintf("observability-agent-%s.dmg", p.version))
	
	// Create temporary directory for DMG contents
	appDir := filepath.Join(p.outDir, "ObservabilityAgent.app")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}

	// Build DMG using hdiutil
	cmd := exec.Command("hdiutil", "create", "-volname", "Observability Agent", "-srcfolder", appDir, "-ov", "-format", "UDZO", dmgFile)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("hdiutil failed: %w", err)
	}

	return dmgFile, nil
}

func (p *Packager) packageHelm() (string, error) {
	// Create Helm chart
	chartDir := filepath.Join(p.outDir, "observability-agent")
	if err := os.MkdirAll(filepath.Join(chartDir, "templates"), 0755); err != nil {
		return "", err
	}

	// Chart.yaml
	chartYaml := fmt.Sprintf(`apiVersion: v2
name: observability-agent
description: Security Observability Agent Helm Chart
type: application
version: %s
appVersion: "%s"
`, p.version, p.version)

	if err := os.WriteFile(filepath.Join(chartDir, "Chart.yaml"), []byte(chartYaml), 0644); err != nil {
		return "", err
	}

	// Package chart
	cmd := exec.Command("helm", "package", chartDir, "-d", p.outDir)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("helm package failed: %w", err)
	}

	// Find the packaged chart
	files, err := filepath.Glob(filepath.Join(p.outDir, "observability-agent-*.tgz"))
	if err != nil || len(files) == 0 {
		return "", fmt.Errorf("no helm chart found")
	}

	return files[0], nil
}

func (p *Packager) packageOVA() (string, error) {
	// OVA packaging requires OVF tool or manual creation
	// This is a placeholder - actual implementation would use ovftool
	ovaFile := filepath.Join(p.outDir, fmt.Sprintf("observability-agent-%s.ova", p.version))
	
	// In a real implementation, this would:
	// 1. Create OVF descriptor
	// 2. Package VMDK disk image
	// 3. Create OVA archive
	
	return ovaFile, nil
}

func (p *Packager) packageDocker() (string, error) {
	// Build Docker image
	imageTag := fmt.Sprintf("observability-agent:%s", p.version)
	cmd := exec.Command("docker", "build", "-t", imageTag, ".")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("docker build failed: %w", err)
	}

	// Save image to tar
	tarFile := filepath.Join(p.outDir, fmt.Sprintf("observability-agent-%s.tar", p.version))
	cmd = exec.Command("docker", "save", "-o", tarFile, imageTag)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("docker save failed: %w", err)
	}

	return tarFile, nil
}

func (p *Packager) packageHomebrew() (string, error) {
	// Create Homebrew formula
	formulaFile := filepath.Join(p.outDir, "observability-agent.rb")
	formula := fmt.Sprintf(`class ObservabilityAgent < Formula
  desc "Security Observability Agent"
  homepage "https://observability.platform"
  url "https://github.com/observability-platform/agent/releases/download/v%s/observability-agent-%s.tar.gz"
  sha256 "YOUR-SHA256-HERE"
  version "%s"

  def install
    bin.install "observability-agent"
  end

  test do
    system "#{bin}/observability-agent", "--version"
  end
end
`, p.version, p.version, p.version)

	if err := os.WriteFile(formulaFile, []byte(formula), 0644); err != nil {
		return "", err
	}

	return formulaFile, nil
}

func (p *Packager) packageChocolatey() (string, error) {
	// Create Chocolatey package
	nuspecFile := filepath.Join(p.outDir, "observability-agent.nuspec")
	nuspec := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<package xmlns="http://schemas.microsoft.com/packaging/2015/06/nuspec.xsd">
  <metadata>
    <id>observability-agent</id>
    <version>%s</version>
    <title>Observability Agent</title>
    <authors>Observability Platform</authors>
    <description>Security Observability Agent</description>
  </metadata>
  <files>
    <file src="observability-agent.exe" target="tools\" />
  </files>
</package>
`, p.version)

	if err := os.WriteFile(nuspecFile, []byte(nuspec), 0644); err != nil {
		return "", err
	}

	// Package Chocolatey (requires choco)
	cmd := exec.Command("choco", "pack", nuspecFile, "--output-directory", p.outDir)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("choco pack failed: %w", err)
	}

	// Find the packaged nupkg
	files, err := filepath.Glob(filepath.Join(p.outDir, "observability-agent.*.nupkg"))
	if err != nil || len(files) == 0 {
		return "", fmt.Errorf("no chocolatey package found")
	}

	return files[0], nil
}
