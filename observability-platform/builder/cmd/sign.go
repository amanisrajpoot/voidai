package cmd

import (
	"fmt"
	"os"

	"github.com/observability-platform/builder/internal/sign"
	"github.com/spf13/cobra"
)

var signKey string
var signCert string
var signPassword string

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign --artifact <path>",
	Short: "Sign an artifact",
	Long: `Sign an artifact for secure distribution.

Supports:
  - macOS: codesign with .p12 certificate
  - Windows: signtool with .pfx certificate
  - Linux: GPG signing for deb/rpm packages

Examples:
  builder sign --artifact ./dist/agent.deb --key ./keys/private.pem
  builder sign --artifact ./dist/agent.dmg --cert ./keys/cert.p12 --password <password>
  builder sign --artifact ./dist/agent.msi --cert ./keys/cert.pfx --password <password>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		artifact, _ := cmd.Flags().GetString("artifact")
		if artifact == "" {
			return fmt.Errorf("--artifact is required")
		}

		if _, err := os.Stat(artifact); os.IsNotExist(err) {
			return fmt.Errorf("artifact not found: %s", artifact)
		}

		signer := sign.NewSigner()
		
		if signKey != "" {
			signer.SetKey(signKey)
		}
		if signCert != "" {
			signer.SetCert(signCert)
		}
		if signPassword != "" {
			signer.SetPassword(signPassword)
		}

		if err := signer.Sign(artifact); err != nil {
			return fmt.Errorf("signing failed: %w", err)
		}

		fmt.Printf("✓ Successfully signed %s\n", artifact)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(signCmd)
	signCmd.Flags().String("artifact", "", "Path to artifact to sign (required)")
	signCmd.Flags().StringVar(&signKey, "key", "", "Path to private key (for GPG/Linux)")
	signCmd.Flags().StringVar(&signCert, "cert", "", "Path to certificate (.p12 for macOS, .pfx for Windows)")
	signCmd.Flags().StringVar(&signPassword, "password", "", "Certificate password (if required)")
}
