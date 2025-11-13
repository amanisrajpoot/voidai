package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign artifacts with cryptographic signatures",
	Long: `Sign artifacts for secure distribution.
Supports: GPG (deb/rpm), codesign (macOS), signtool (Windows)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		artifact, _ := cmd.Flags().GetString("artifact")
		key, _ := cmd.Flags().GetString("key")
		method, _ := cmd.Flags().GetString("method")

		if artifact == "" {
			return fmt.Errorf("--artifact is required")
		}

		if method == "" {
			// Auto-detect from artifact extension
			ext := filepath.Ext(artifact)
			switch ext {
			case ".deb", ".rpm":
				method = "gpg"
			case ".dmg", ".pkg":
				method = "codesign"
			case ".msi", ".exe":
				method = "signtool"
			default:
				return fmt.Errorf("cannot auto-detect signing method, use --method")
			}
		}

		return signArtifact(artifact, key, method)
	},
}

func init() {
	signCmd.Flags().String("artifact", "", "Artifact to sign")
	signCmd.Flags().String("key", "", "Signing key path")
	signCmd.Flags().String("method", "", "Signing method (gpg, codesign, signtool)")
	rootCmd.AddCommand(signCmd)
}

func signArtifact(artifact, key, method string) error {
	fmt.Printf("Signing %s using %s...\n", artifact, method)

	switch method {
	case "gpg":
		return signGPG(artifact, key)
	case "codesign":
		return signCodesign(artifact, key)
	case "signtool":
		return signSigntool(artifact, key)
	default:
		return fmt.Errorf("unsupported signing method: %s", method)
	}
}

func signGPG(artifact, key string) error {
	// GPG signing for deb/rpm
	var cmd *exec.Cmd
	if key != "" {
		cmd = exec.Command("gpg", "--default-key", key, "--armor", "--detach-sign", artifact)
	} else {
		cmd = exec.Command("gpg", "--armor", "--detach-sign", artifact)
	}
	return cmd.Run()
}

func signCodesign(artifact, key string) error {
	// macOS codesign
	if key == "" {
		key = "-" // Use default identity
	}
	cmd := exec.Command("codesign", "--sign", key, "--timestamp", artifact)
	return cmd.Run()
}

func signSigntool(artifact, key string) error {
	// Windows signtool
	if key == "" {
		return fmt.Errorf("signtool requires --key (certificate path)")
	}
	cmd := exec.Command("signtool", "sign", "/f", key, "/t", "http://timestamp.digicert.com", artifact)
	return cmd.Run()
}
