package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/observability-platform/builder/internal/upload"
	"github.com/spf13/cobra"
)

var uploadRepo string
var uploadRepoType string
var uploadToken string

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload --artifact <path>",
	Short: "Upload artifacts to repository",
	Long: `Upload signed artifacts to package repositories.

Supported repositories:
  - artifactory (JFrog Artifactory)
  - npm (npm registry)
  - pypi (PyPI)
  - maven (Maven Central)
  - nuget (NuGet Gallery)
  - docker (Docker Hub/Registry)
  - github (GitHub Releases)
  - s3 (AWS S3)

Examples:
  builder upload --artifact ./dist/*.deb --repo artifactory --repo-type=artifactory --token <token>
  builder upload --artifact ./dist/agent-*.tgz --repo npm --repo-type=npm --token <token>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		artifactPattern, _ := cmd.Flags().GetString("artifact")
		if artifactPattern == "" {
			return fmt.Errorf("--artifact is required")
		}

		if uploadRepo == "" {
			return fmt.Errorf("--repo is required")
		}

		if uploadRepoType == "" {
			uploadRepoType = "artifactory" // default
		}

		matches, err := filepath.Glob(artifactPattern)
		if err != nil {
			return fmt.Errorf("invalid artifact pattern: %w", err)
		}

		if len(matches) == 0 {
			return fmt.Errorf("no artifacts found matching pattern: %s", artifactPattern)
		}

		uploader := upload.NewUploader(uploadRepo, uploadRepoType)
		if uploadToken != "" {
			uploader.SetToken(uploadToken)
		}

		for _, artifact := range matches {
			if err := uploader.Upload(artifact); err != nil {
				return fmt.Errorf("failed to upload %s: %w", artifact, err)
			}
			fmt.Printf("✓ Uploaded %s\n", artifact)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().String("artifact", "", "Path/pattern to artifact(s) to upload (required)")
	uploadCmd.Flags().StringVar(&uploadRepo, "repo", "", "Repository URL or identifier (required)")
	uploadCmd.Flags().StringVar(&uploadRepoType, "repo-type", "artifactory", "Repository type")
	uploadCmd.Flags().StringVar(&uploadToken, "token", "", "Authentication token")
}
