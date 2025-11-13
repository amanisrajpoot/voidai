package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload artifacts to repositories",
	Long: `Upload signed artifacts to package repositories.
Supports: npm, pypi, maven, nuget, artifactory, s3`,
	RunE: func(cmd *cobra.Command, args []string) error {
		artifact, _ := cmd.Flags().GetString("artifact")
		repo, _ := cmd.Flags().GetString("repo")
		repoType, _ := cmd.Flags().GetString("type")

		if artifact == "" {
			return fmt.Errorf("--artifact is required (supports glob patterns)")
		}

		if repo == "" {
			return fmt.Errorf("--repo is required")
		}

		if repoType == "" {
			// Auto-detect from artifact extension
			ext := filepath.Ext(artifact)
			switch ext {
			case ".tgz", ".tar.gz":
				repoType = "npm"
			case ".whl", ".tar.gz":
				repoType = "pypi"
			case ".jar":
				repoType = "maven"
			case ".nupkg":
				repoType = "nuget"
			case ".deb", ".rpm":
				repoType = "apt" // or "yum"
			default:
				return fmt.Errorf("cannot auto-detect repo type, use --type")
			}
		}

		// Expand glob patterns
		matches, err := filepath.Glob(artifact)
		if err != nil {
			return err
		}

		for _, match := range matches {
			if err := uploadArtifact(match, repo, repoType); err != nil {
				return fmt.Errorf("failed to upload %s: %w", match, err)
			}
		}

		return nil
	},
}

func init() {
	uploadCmd.Flags().String("artifact", "", "Artifact(s) to upload (supports glob)")
	uploadCmd.Flags().String("repo", "", "Repository URL or name")
	uploadCmd.Flags().String("type", "", "Repository type (npm, pypi, maven, nuget, apt, yum, s3)")
	rootCmd.AddCommand(uploadCmd)
}

func uploadArtifact(artifact, repo, repoType string) error {
	fmt.Printf("Uploading %s to %s (%s)...\n", artifact, repo, repoType)

	switch repoType {
	case "npm":
		return uploadNPM(artifact, repo)
	case "pypi":
		return uploadPyPI(artifact, repo)
	case "maven":
		return uploadMaven(artifact, repo)
	case "nuget":
		return uploadNuGet(artifact, repo)
	case "apt", "yum":
		return uploadPackageRepo(artifact, repo, repoType)
	case "s3":
		return uploadS3(artifact, repo)
	case "artifactory":
		return uploadArtifactory(artifact, repo)
	default:
		return fmt.Errorf("unsupported repository type: %s", repoType)
	}
}

func uploadNPM(artifact, repo string) error {
	// npm publish
	cmd := exec.Command("npm", "publish", artifact)
	if repo != "npmjs.com" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("NPM_REGISTRY=%s", repo))
	}
	return cmd.Run()
}

func uploadPyPI(artifact, repo string) error {
	// twine upload
	cmd := exec.Command("twine", "upload", "--repository-url", repo, artifact)
	return cmd.Run()
}

func uploadMaven(artifact, repo string) error {
	// mvn deploy
	cmd := exec.Command("mvn", "deploy:deploy-file", "-Dfile="+artifact, "-DrepositoryId="+repo)
	return cmd.Run()
}

func uploadNuGet(artifact, repo string) error {
	// nuget push
	cmd := exec.Command("nuget", "push", artifact, "-Source", repo)
	return cmd.Run()
}

func uploadPackageRepo(artifact, repo, repoType string) error {
	if repoType == "apt" {
		// reprepro or aptly
		cmd := exec.Command("reprepro", "-b", repo, "includedeb", "stable", artifact)
		return cmd.Run()
	} else {
		// createrepo for yum
		cmd := exec.Command("createrepo", repo)
		return cmd.Run()
	}
}

func uploadS3(artifact, repo string) error {
	// aws s3 cp
	key := filepath.Base(artifact)
	if !strings.HasSuffix(repo, "/") {
		repo += "/"
	}
	cmd := exec.Command("aws", "s3", "cp", artifact, repo+key)
	return cmd.Run()
}

func uploadArtifactory(artifact, repo string) error {
	// jfrog or curl
	key := filepath.Base(artifact)
	url := fmt.Sprintf("%s/%s", repo, key)
	cmd := exec.Command("curl", "-T", artifact, "-u", "$ARTIFACTORY_USER:$ARTIFACTORY_PASSWORD", url)
	return cmd.Run()
}
