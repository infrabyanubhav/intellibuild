package ci

import (
	"fmt"
	"os"
	"os/exec"
)

// BuildProject builds the project based on its language (Go, Makefile, Python, npm).
func BuildProject(dir string) error {
	// Check for Go project (go.mod)
	if _, err := os.Stat(fmt.Sprintf("%s/go.mod", dir)); err == nil {
		cmd := exec.Command("go", "build", "./...")
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// Check for Makefile
	if _, err := os.Stat(fmt.Sprintf("%s/Makefile", dir)); err == nil {
		cmd := exec.Command("make")
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// Check for Python project (requirements.txt or setup.py)
	// Check for Python project (requirements.txt or setup.py)
	if _, err := os.Stat(fmt.Sprintf("%s/requirements.txt", dir)); err == nil {
		cmd := exec.Command("python", "setup.py", "install") // For example, build with setup.py
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	} else if _, err := os.Stat(fmt.Sprintf("%s/setup.py", dir)); err == nil {
		cmd := exec.Command("python", "setup.py", "install") // For example, build with setup.py
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// Check for npm project (package.json)
	if _, err := os.Stat(fmt.Sprintf("%s/package.json", dir)); err == nil {
		cmd := exec.Command("npm", "install") // Install dependencies (npm build step)
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// If no supported build command is found, return an error
	return fmt.Errorf("unsupported language or build system")
}
