package sign

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Signer struct {
	key      string
	cert     string
	password string
}

func NewSigner() *Signer {
	return &Signer{}
}

func (s *Signer) SetKey(key string) {
	s.key = key
}

func (s *Signer) SetCert(cert string) {
	s.cert = cert
}

func (s *Signer) SetPassword(password string) {
	s.password = password
}

func (s *Signer) Sign(artifact string) error {
	ext := strings.ToLower(filepath.Ext(artifact))
	
	switch ext {
	case ".deb":
		return s.signDeb(artifact)
	case ".rpm":
		return s.signRPM(artifact)
	case ".dmg", ".pkg":
		return s.signMacOS(artifact)
	case ".msi", ".exe":
		return s.signWindows(artifact)
	default:
		return fmt.Errorf("unsupported artifact type for signing: %s", ext)
	}
}

func (s *Signer) signDeb(artifact string) error {
	// Sign DEB with GPG
	if s.key == "" {
		return fmt.Errorf("GPG key required for DEB signing")
	}

	cmd := exec.Command("dpkg-sig", "--sign", "builder", "-k", s.key, artifact)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (s *Signer) signRPM(artifact string) error {
	// Sign RPM with GPG
	if s.key == "" {
		return fmt.Errorf("GPG key required for RPM signing")
	}

	cmd := exec.Command("rpm", "--addsign", artifact)
	cmd.Env = append(os.Environ(), fmt.Sprintf("GPG_TTY=$(tty)"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (s *Signer) signMacOS(artifact string) error {
	// Sign macOS artifacts with codesign
	if s.cert == "" {
		return fmt.Errorf("certificate (.p12) required for macOS signing")
	}

	// Import certificate to keychain (if needed)
	// codesign requires certificate to be in keychain
	
	// Sign the artifact
	cmd := exec.Command("codesign", "--sign", "Developer ID Application", "--timestamp", artifact)
	if s.password != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("KEYCHAIN_PASSWORD=%s", s.password))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("codesign failed: %w", err)
	}

	// Notarize (optional, requires Apple Developer account)
	// This would typically be done separately via notarytool or altool

	return nil
}

func (s *Signer) signWindows(artifact string) error {
	// Sign Windows artifacts with signtool
	if s.cert == "" {
		return fmt.Errorf("certificate (.pfx) required for Windows signing")
	}

	passwordFlag := ""
	if s.password != "" {
		passwordFlag = fmt.Sprintf("/p:%s", s.password)
	}

	// Find signtool (usually in Windows SDK)
	signtoolPath := "signtool.exe"
	if path := os.Getenv("SIGNTOOL_PATH"); path != "" {
		signtoolPath = path
	}

	cmd := exec.Command(signtoolPath, "sign", "/f", s.cert, passwordFlag, "/t", "http://timestamp.digicert.com", artifact)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}
