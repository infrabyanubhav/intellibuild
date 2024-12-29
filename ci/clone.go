package ci

import (
	"os"
	"os/exec"
)

// CloneRepository clones the given Git repository URL into the specified directory.
func CloneRepository(repoURL, dir string) error {
	cmd := exec.Command("git", "clone", repoURL, dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
