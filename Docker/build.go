package docker

import (
	"fmt"
	"os"
	"os/exec"
)

// BuildDockerImage builds a Docker image from the provided directory.
func BuildDockerImage(dir string) error {
	if _, err := os.Stat(fmt.Sprintf("%s/Dockerfile", dir)); err != nil {
		return fmt.Errorf("Dockerfile not found")
	}

	imageName := "myapp:latest" // You can modify this for dynamic tagging
	cmd := exec.Command("docker", "build", "-t", imageName, dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build Docker image: %w", err)
	}

	return nil
}
