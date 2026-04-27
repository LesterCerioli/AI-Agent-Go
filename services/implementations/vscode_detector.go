// services/implementations/vscode_detector.go
package implementations

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type VSCodeDetector struct {
	workspacePath string
}

type WorkspaceInfo struct {
	Path  string   `json:"path"`
	Name  string   `json:"name"`
	Files []string `json:"files"`
}

func NewVSCodeDetector() (*VSCodeDetector, error) {
	return &VSCodeDetector{}, nil
}

func (d *VSCodeDetector) GetCurrentWorkspace() (string, error) {
	log.Printf("[INFO] Detecting current VS Code workspace")

	// Method 1: Get from current working directory
	if path, err := d.getCurrentDirectory(); err == nil && path != "" {
		log.Printf("[INFO] Using current directory: %s", path)
		return path, nil
	}

	// Method 2: Try to find from VS Code process (Windows)
	if runtime.GOOS == "windows" {
		if path, err := d.getVSCodeWorkspaceWindows(); err == nil && path != "" {
			log.Printf("[INFO] Detected VS Code workspace: %s", path)
			return path, nil
		}
	}

	// Method 3: Use home directory as fallback
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	projectsDir := filepath.Join(homeDir, "projects")
	if err := os.MkdirAll(projectsDir, 0755); err != nil {
		return homeDir, nil
	}

	log.Printf("[INFO] Using fallback directory: %s", projectsDir)
	return projectsDir, nil
}

func (d *VSCodeDetector) getCurrentDirectory() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Check if current directory contains project indicators
	indicators := []string{"go.mod", "package.json", ".git", "requirements.txt", "Cargo.toml"}
	for _, indicator := range indicators {
		if d.fileExists(filepath.Join(currentDir, indicator)) {
			return currentDir, nil
		}
	}

	return currentDir, nil
}

func (d *VSCodeDetector) getVSCodeWorkspaceWindows() (string, error) {
	// Get VS Code process
	cmd := exec.Command("powershell", "-Command",
		"Get-Process -Name Code -ErrorAction SilentlyContinue | Select-Object -First 1 -ExpandProperty Id")

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	processID := strings.TrimSpace(string(output))
	if processID == "" {
		return "", fmt.Errorf("VS Code not running")
	}

	// Get working directory from process
	cmd = exec.Command("powershell", "-Command",
		fmt.Sprintf("(Get-Process -Id %s).StartInfo.WorkingDirectory", processID))

	output, err = cmd.Output()
	if err != nil {
		return "", err
	}

	workDir := strings.TrimSpace(string(output))
	if workDir != "" && d.fileExists(workDir) {
		return workDir, nil
	}

	return "", fmt.Errorf("could not detect workspace")
}

func (d *VSCodeDetector) fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (d *VSCodeDetector) ScanProjectFiles(projectPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip hidden directories and node_modules
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" || name == "__pycache__" {
				return filepath.SkipDir
			}
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(projectPath, path)
		if err == nil {
			files = append(files, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan project: %w", err)
	}

	return files, nil
}

func (d *VSCodeDetector) DetectProjectLanguage(projectPath string) string {
	if d.fileExists(filepath.Join(projectPath, "go.mod")) {
		return "go"
	}
	if d.fileExists(filepath.Join(projectPath, "package.json")) {
		return "javascript"
	}
	if d.fileExists(filepath.Join(projectPath, "requirements.txt")) {
		return "python"
	}
	if d.fileExists(filepath.Join(projectPath, "Cargo.toml")) {
		return "rust"
	}
	if d.fileExists(filepath.Join(projectPath, "pubspec.yaml")) {
		return "dart"
	}
	return "unknown"
}
