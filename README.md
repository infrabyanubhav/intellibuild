# IntelliBuildCI

A lightweight, self-hosted Continuous Integration (CI) webhook server written in Go that automates build, test, and security scanning processes.

## Overview

IntelliBuildCI is a custom CI/CD pipeline orchestrator that:
- Receives webhook events from Git providers (GitHub, GitLab, Bitbucket, etc.)
- Clones repositories automatically
- Builds projects (supports Go, Makefile, Python, npm/Node.js)
- Runs tests
- Builds Docker images
- Performs security scanning with Trivy

## Features

- **Multi-language support**: Automatically detects and builds Go, Python, Node.js, and Makefile projects
- **Docker integration**: Builds Docker images from Dockerfiles
- **Security scanning**: Integrates Trivy for container vulnerability scanning
- **Webhook-driven**: Triggered by Git webhook events

## Architecture

```
IntelliBuildCI/
├── main.go              # HTTP server entry point
├── webhook/             # Webhook request handling
│   └── webhook.go       # Main webhook handler and pipeline orchestration
├── ci/                  # CI operations
│   ├── clone.go         # Git repository cloning
│   ├── build.go         # Project building (multi-language)
│   └── test.go          # Test execution (multi-language)
└── docker/              # Docker operations
    ├── build.go         # Docker image building
    └── scan.go          # Trivy security scanning
```

## Installation

### Prerequisites

- Go 1.23.2 or higher
- Docker (for building images)
- Trivy (for security scanning)
- Git

### Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd intellibuildCI
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o intellibuildci
```

4. Run the server:
```bash
./intellibuildci
```

The server will start on port `8080` and listen for webhook requests at `/webhook`.

## Usage

Send a POST request to `http://localhost:8080/webhook` with a JSON payload containing the repository URL:

```json
{
  "repoURL": "https://github.com/user/repo.git"
}
```

The pipeline will:
1. Clone the repository
2. Build the project
3. Run tests
4. Build a Docker image (if Dockerfile exists)
5. Scan the Docker image with Trivy

## Security Issues

⚠️ **WARNING: This project contains several critical security vulnerabilities. Do not use in production without addressing these issues.**

### Critical Issues

#### 1. Hardcoded Repository URL
**Location**: `webhook/webhook.go:80`

The `parseRepoURL` function contains a hardcoded fallback URL:
```go
return "https://github.com/user/repo.git" // Modify accordingly
```

**Risk**: If payload parsing fails, the system defaults to a hardcoded repository, which could lead to:
- Unintended builds from wrong repositories
- Potential supply chain attacks
- Unauthorized code execution

**Recommendation**: 
- Remove the hardcoded fallback
- Return an error if the repository URL cannot be parsed
- Implement proper JSON parsing with validation

#### 2. No Webhook Signature Validation
**Location**: `webhook/webhook.go:17-21`

The webhook handler accepts any POST request without verifying the source:
```go
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
```

**Risk**: 
- Anyone who knows the webhook URL can trigger the CI/CD pipeline
- Potential for:
  - Denial of Service (DoS) attacks
  - Resource exhaustion
  - Unauthorized code execution
  - Malicious code injection

**Recommendation**:
- Implement webhook signature validation for GitHub (X-Hub-Signature-256)
- Implement webhook signature validation for GitLab (X-Gitlab-Token)
- Implement webhook signature validation for Bitbucket
- Verify signatures using HMAC-SHA256

#### 3. No Authentication/Authorization
**Location**: `main.go:12`

The webhook endpoint is completely open:
```go
http.HandleFunc("/webhook", webhook.WebhookHandler)
```

**Risk**:
- Public endpoint accessible to anyone
- No access control or rate limiting
- Vulnerable to brute force attacks

**Recommendation**:
- Implement API key authentication
- Add bearer token authentication
- Implement basic authentication
- Add IP whitelisting for production use
- Implement rate limiting to prevent abuse

### High Priority Issues

#### 4. Missing Credential Handling for Git Operations
**Location**: `ci/clone.go:9-14`

Git cloning doesn't support authentication:
```go
func CloneRepository(repoURL, dir string) error {
	cmd := exec.Command("git", "clone", repoURL, dir)
	// ...
}
```

**Risk**:
- Cannot clone private repositories
- If credentials are needed, they would need to be passed insecurely
- No support for SSH keys or token-based authentication

**Recommendation**:
- Support SSH key-based authentication
- Support token-based authentication (store tokens securely)
- Use environment variables or secrets manager for credentials
- Never hardcode credentials in the codebase
- Support git credential helpers

#### 5. Potential Command Injection
**Location**: `ci/clone.go:10`, `webhook/webhook.go:32`

Repository URLs are used directly in shell commands without sanitization:
```go
repoURL := parseRepoURL(payload)
cmd := exec.Command("git", "clone", repoURL, dir)
```

**Risk**:
- If `repoURL` contains malicious characters, it could lead to command injection
- Attackers could execute arbitrary commands on the server

**Recommendation**:
- Validate and sanitize repository URLs
- Use whitelist-based validation for URL patterns
- Ensure URLs match expected format (e.g., `https://github.com/...` or `git@github.com:...`)
- Use `exec.Command` with separate arguments (already done, but validate inputs)

### Medium Priority Issues

#### 6. Information Disclosure in Error Messages
**Location**: `webhook/webhook.go:42-43`, `webhook/webhook.go:48-49`

Error messages expose internal details:
```go
http.Error(w, fmt.Sprintf("Error cloning repository: %v", err), http.StatusInternalServerError)
```

**Risk**:
- Error messages may leak:
  - File system paths
  - Internal system information
  - Stack traces
  - Configuration details

**Recommendation**:
- Sanitize error messages before returning to clients
- Log detailed errors server-side
- Return generic error messages to clients
- Implement structured logging

#### 7. No Input Validation
**Location**: `webhook/webhook.go:24-31`

The webhook payload is read without validation:
```go
body, err := io.ReadAll(r.Body)
payload := string(body)
repoURL := parseRepoURL(payload)
```

**Risk**:
- No size limits on request body (could lead to DoS)
- No content-type validation
- No JSON schema validation

**Recommendation**:
- Limit request body size
- Validate Content-Type header
- Implement JSON schema validation
- Add timeout handling

#### 8. Insecure Temporary Directory Usage
**Location**: `webhook/webhook.go:40`

Temporary directories are created in `/tmp`:
```go
dir := fmt.Sprintf("/tmp/repo-%d", time.Now().Unix())
```

**Risk**:
- Potential race conditions if multiple requests arrive simultaneously
- Temporary files not cleaned up on errors
- `/tmp` directory may have insecure permissions

**Recommendation**:
- Use `os.MkdirTemp` for secure temporary directory creation
- Implement cleanup of temporary directories
- Ensure proper file permissions
- Add cleanup on panic/error

## Security Best Practices to Implement

1. **Webhook Security**:
   - Validate webhook signatures from all Git providers
   - Implement request rate limiting
   - Add request timeout handling

2. **Authentication**:
   - Add API key or token-based authentication
   - Implement IP whitelisting for production
   - Use HTTPS only

3. **Credential Management**:
   - Use environment variables for sensitive data
   - Integrate with secrets managers (AWS Secrets Manager, HashiCorp Vault, etc.)
   - Never commit credentials to version control

4. **Input Validation**:
   - Validate all inputs (URLs, payloads, etc.)
   - Implement size limits
   - Use allowlists instead of blocklists

5. **Error Handling**:
   - Sanitize error messages
   - Implement structured logging
   - Don't expose internal details to clients

6. **Resource Management**:
   - Clean up temporary directories
   - Implement resource limits
   - Add timeout handling for long-running operations

7. **Monitoring**:
   - Add logging and monitoring
   - Implement alerting for suspicious activities
   - Track webhook request patterns

## Contributing

Contributions are welcome! Please ensure:
- All security issues are addressed
- Code follows Go best practices
- Tests are included for new features
- Documentation is updated

## License

[Add your license here]

## Disclaimer

This project is provided as-is for educational purposes. **Do not use in production environments without addressing all security vulnerabilities mentioned above.**

