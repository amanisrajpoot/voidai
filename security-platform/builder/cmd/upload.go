package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var uploadArtifact string
var uploadRepo string
var uploadRepoType string

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload artifacts to repository",
	Long: `Upload signed artifacts to package repositories.
Supports: artifactory, npm, pypi, maven, nuget, docker`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if uploadArtifact == "" {
			return fmt.Errorf("--artifact is required")
		}

		if uploadRepo == "" {
			return fmt.Errorf("--repo is required")
		}

		// Determine repo type from artifact extension if not specified
		if uploadRepoType == "" {
			uploadRepoType = detectRepoType(uploadArtifact)
		}

		fmt.Printf("Uploading %s to %s (%s)...\n", uploadArtifact, uploadRepo, uploadRepoType)

		switch uploadRepoType {
		case "artifactory":
			return uploadToArtifactory(uploadArtifact, uploadRepo)
		case "npm":
			return uploadToNPM(uploadArtifact, uploadRepo)
		case "pypi":
			return uploadToPyPI(uploadArtifact, uploadRepo)
		case "maven":
			return uploadToMaven(uploadArtifact, uploadRepo)
		case "nuget":
			return uploadToNuGet(uploadArtifact, uploadRepo)
		case "docker":
			return uploadToDocker(uploadArtifact, uploadRepo)
		default:
			return fmt.Errorf("unsupported repository type: %s", uploadRepoType)
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringVar(&uploadArtifact, "artifact", "", "artifact file or glob pattern")
	uploadCmd.Flags().StringVar(&uploadRepo, "repo", "", "repository URL or name")
	uploadCmd.Flags().StringVar(&uploadRepoType, "type", "", "repository type (artifactory, npm, pypi, maven, nuget, docker)")
}

func detectRepoType(artifact string) string {
	ext := strings.ToLower(getFileExtension(artifact))
	switch ext {
	case ".tgz", ".tar.gz":
		if strings.Contains(artifact, "package") {
			return "npm"
		}
		return "artifactory"
	case ".whl", ".tar.gz":
		if strings.Contains(artifact, "dist") {
			return "pypi"
		}
		return "artifactory"
	case ".jar", ".pom":
		return "maven"
	case ".nupkg":
		return "nuget"
	case ".deb", ".rpm", ".msi", ".dmg":
		return "artifactory"
	default:
		return "artifactory"
	}
}

func getFileExtension(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i:]
		}
	}
	return ""
}

func uploadToArtifactory(artifact, repo string) error {
	fmt.Println("  Uploading to Artifactory")
	// In real implementation:
	// curl -X PUT -T artifact -u user:pass repo/artifact
	return nil
}

func uploadToNPM(artifact, repo string) error {
	fmt.Println("  Publishing to NPM")
	// In real implementation:
	// npm publish --registry=repo
	return nil
}

func uploadToPyPI(artifact, repo string) error {
	fmt.Println("  Uploading to PyPI")
	// In real implementation:
	// twine upload --repository-url=repo artifact
	return nil
}

func uploadToMaven(artifact, repo string) error {
	fmt.Println("  Deploying to Maven")
	// In real implementation:
	// mvn deploy:deploy-file -Dfile=artifact -DrepositoryId=repo
	return nil
}

func uploadToNuGet(artifact, repo string) error {
	fmt.Println("  Pushing to NuGet")
	// In real implementation:
	// dotnet nuget push artifact --source=repo
	return nil
}

func uploadToDocker(artifact, repo string) error {
	fmt.Println("  Pushing to Docker registry")
	// In real implementation:
	// docker push repo/image:tag
	return nil
}
