package ci

import (
	"fmt"
	"os"
	"os/exec"
)

// RunTests runs the tests based on the detected language (Go, Makefile, Python, npm).
func RunTests(dir string) error {
	// Check for Go project (go.mod)
	if _, err := os.Stat(fmt.Sprintf("%s/go.mod", dir)); err == nil {
		cmd := exec.Command("go", "test", "./...")
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// Check for Makefile
	if _, err := os.Stat(fmt.Sprintf("%s/Makefile", dir)); err == nil {
		cmd := exec.Command("make", "test")
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// Check for Python project (requirements.txt)
	if _, err := os.Stat(fmt.Sprintf("%s/requirements.txt", dir)); err == nil {
		cmd := exec.Command("pytest") // assuming pytest is installed
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// Check for npm project (package.json)
	if _, err := os.Stat(fmt.Sprintf("%s/package.json", dir)); err == nil {
		cmd := exec.Command("npm", "test")
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// If no test command is found, return an error
	return fmt.Errorf("no test command available for this project")
}
