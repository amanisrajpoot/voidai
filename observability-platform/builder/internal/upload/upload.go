package upload

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Uploader struct {
	repo     string
	repoType string
	token    string
}

func NewUploader(repo, repoType string) *Uploader {
	return &Uploader{
		repo:     repo,
		repoType: repoType,
	}
}

func (u *Uploader) SetToken(token string) {
	u.token = token
}

func (u *Uploader) Upload(artifact string) error {
	switch u.repoType {
	case "artifactory":
		return u.uploadToArtifactory(artifact)
	case "npm":
		return u.uploadToNPM(artifact)
	case "pypi":
		return u.uploadToPyPI(artifact)
	case "maven":
		return u.uploadToMaven(artifact)
	case "nuget":
		return u.uploadToNuGet(artifact)
	case "docker":
		return u.uploadToDocker(artifact)
	case "github":
		return u.uploadToGitHub(artifact)
	case "s3":
		return u.uploadToS3(artifact)
	default:
		return fmt.Errorf("unsupported repository type: %s", u.repoType)
	}
}

func (u *Uploader) uploadToArtifactory(artifact string) error {
	// Upload to JFrog Artifactory
	filename := filepath.Base(artifact)
	url := fmt.Sprintf("%s/%s", strings.TrimSuffix(u.repo, "/"), filename)

	file, err := os.Open(artifact)
	if err != nil {
		return err
	}
	defer file.Close()

	req, err := http.NewRequest("PUT", url, file)
	if err != nil {
		return err
	}

	if u.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", u.token))
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed: %s - %s", resp.Status, string(body))
	}

	return nil
}

func (u *Uploader) uploadToNPM(artifact string) error {
	// Upload to npm registry using npm publish
	cmd := exec.Command("npm", "publish", artifact)
	if u.token != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("NPM_TOKEN=%s", u.token))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (u *Uploader) uploadToPyPI(artifact string) error {
	// Upload to PyPI using twine
	cmd := exec.Command("twine", "upload", artifact)
	if u.token != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("TWINE_PASSWORD=%s", u.token))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (u *Uploader) uploadToMaven(artifact string) error {
	// Upload to Maven Central using mvn deploy
	cmd := exec.Command("mvn", "deploy:deploy-file", 
		"-Dfile="+artifact,
		"-DrepositoryId="+u.repo,
		"-Durl="+u.repo,
	)
	if u.token != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("MAVEN_PASSWORD=%s", u.token))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (u *Uploader) uploadToNuGet(artifact string) error {
	// Upload to NuGet using nuget push
	cmd := exec.Command("nuget", "push", artifact, "-Source", u.repo, "-ApiKey", u.token)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (u *Uploader) uploadToDocker(artifact string) error {
	// Load and push Docker image
	cmd := exec.Command("docker", "load", "-i", artifact)
	if err := cmd.Run(); err != nil {
		return err
	}

	// Extract image name from artifact or use repo
	imageName := u.repo
	cmd = exec.Command("docker", "push", imageName)
	if u.token != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("DOCKER_PASSWORD=%s", u.token))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (u *Uploader) uploadToGitHub(artifact string) error {
	// Upload to GitHub Releases (requires gh CLI or API)
	// This is a simplified version - real implementation would use GitHub API
	cmd := exec.Command("gh", "release", "upload", "latest", artifact, "--repo", u.repo)
	if u.token != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("GITHUB_TOKEN=%s", u.token))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (u *Uploader) uploadToS3(artifact string) error {
	// Upload to S3 using AWS CLI
	filename := filepath.Base(artifact)
	s3Path := fmt.Sprintf("%s/%s", strings.TrimSuffix(u.repo, "/"), filename)
	
	cmd := exec.Command("aws", "s3", "cp", artifact, s3Path)
	if u.token != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", u.token))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
