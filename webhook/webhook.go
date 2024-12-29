package webhook

import (
	"fmt"
	"io"

	"net/http"
	"strings"
	"time"

	// Corrected import path for CI
	"IntelliBuildCI/v/ci"
	"IntelliBuildCI/v/docker" // Corrected import path for Docker
)

// WebhookHandler handles incoming Git webhook events.
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read and log the request body (payload)
	body, err := io.ReadAll(r.Body) // Replacing ioutil.ReadAll with io.ReadAll
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Assuming the payload contains a `repoURL` key
	payload := string(body)
	repoURL := parseRepoURL(payload) // You will need a function to extract the repo URL

	if repoURL == "" {
		http.Error(w, "Repository URL not found", http.StatusBadRequest)
		return
	}

	// Step 1: Clone the repo
	dir := fmt.Sprintf("/tmp/repo-%d", time.Now().Unix())
	if err := ci.CloneRepository(repoURL, dir); err != nil {
		http.Error(w, fmt.Sprintf("Error cloning repository: %v", err), http.StatusInternalServerError)
		return
	}

	// Step 2: Build the project
	if err := ci.BuildProject(dir); err != nil {
		http.Error(w, fmt.Sprintf("Build failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Step 3: Run tests (optional based on the webhook payload or config)
	if err := ci.RunTests(dir); err != nil {
		http.Error(w, fmt.Sprintf("Tests failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Step 4: Build Docker image
	if err := docker.BuildDockerImage(dir); err != nil {
		http.Error(w, fmt.Sprintf("Docker build failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Step 5: Run Trivy scan
	if err := docker.RunTrivyScan("myapp:latest"); err != nil {
		http.Error(w, fmt.Sprintf("Trivy scan failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Responding with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("CI/CD pipeline executed successfully"))
}

// Utility function to parse the repository URL from the webhook payload (example).
func parseRepoURL(payload string) string {
	// Assuming the payload contains a URL in a certain format (e.g., JSON).
	if strings.Contains(payload, "repoURL") {
		// Example for simple payload parsing. Use a real parser like JSON or XML as needed.
		return "https://github.com/user/repo.git" // Modify accordingly
	}
	return ""
}
