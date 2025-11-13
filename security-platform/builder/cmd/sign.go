package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var signArtifact string
var signKey string

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign artifact with cryptographic signature",
	Long: `Sign an artifact (package, binary, etc.) using the specified key.
Supports GPG for Linux packages, codesign for macOS, and signtool for Windows.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if signArtifact == "" {
			return fmt.Errorf("--artifact is required")
		}

		if signKey == "" {
			return fmt.Errorf("--key is required")
		}

		if _, err := os.Stat(signArtifact); os.IsNotExist(err) {
			return fmt.Errorf("artifact not found: %s", signArtifact)
		}

		if _, err := os.Stat(signKey); os.IsNotExist(err) {
			return fmt.Errorf("key file not found: %s", signKey)
		}

		fmt.Printf("Signing %s with key %s...\n", signArtifact, signKey)

		// Determine signing method based on artifact extension
		ext := getFileExtension(signArtifact)
		switch ext {
		case ".deb", ".rpm":
			return signWithGPG(signArtifact, signKey)
		case ".dmg", ".pkg":
			return signWithCodesign(signArtifact, signKey)
		case ".msi", ".exe":
			return signWithSigntool(signArtifact, signKey)
		default:
			return fmt.Errorf("unsupported artifact type: %s", ext)
		}
	},
}

func init() {
	rootCmd.AddCommand(signCmd)
	signCmd.Flags().StringVar(&signArtifact, "artifact", "", "artifact file to sign")
	signCmd.Flags().StringVar(&signKey, "key", "", "private key file or certificate")
}

func getFileExtension(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i:]
		}
	}
	return ""
}

func signWithGPG(artifact, key string) error {
	fmt.Println("  Using GPG signing for Linux package")
	// In real implementation:
	// gpg --import key
	// dpkg-sig --sign builder -k <key-id> artifact.deb
	// or rpm --addsign artifact.rpm
	return nil
}

func signWithCodesign(artifact, key string) error {
	fmt.Println("  Using codesign for macOS package")
	// In real implementation:
	// codesign --sign "Developer ID Application: ..." --timestamp artifact.dmg
	return nil
}

func signWithSigntool(artifact, key string) error {
	fmt.Println("  Using signtool for Windows package")
	// In real implementation:
	// signtool sign /f key.pfx /p password artifact.msi
	return nil
}
