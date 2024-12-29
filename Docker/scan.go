package docker

import (
	"fmt"
	"os"
	"os/exec"
)

// RunTrivyScan runs a security scan using Trivy on the given Docker image.
func RunTrivyScan(imageName string) error {
	cmd := exec.Command("trivy", "image", "--no-progress", imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Trivy scan failed: %w", err)
	}
	return nil
}
