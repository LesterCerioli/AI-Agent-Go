// services/implementations/file_generator.go
package implementations

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type FileGenerator struct{}

func NewFileGenerator() (*FileGenerator, error) {
	return &FileGenerator{}, nil
}

func (g *FileGenerator) ParseAndGenerateFiles(projectPath string, deepseekResponse string) ([]string, error) {
	log.Printf("[INFO] Parsing DeepSeek response and generating files in: %s", projectPath)

	var files map[string]string

	// Try to parse as JSON
	if err := json.Unmarshal([]byte(deepseekResponse), &files); err != nil {
		// If not JSON, try to extract code blocks
		files = g.extractCodeBlocks(deepseekResponse)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files found in DeepSeek response")
	}

	var createdFiles []string

	for filePath, content := range files {
		// Clean up file path
		cleanPath := strings.TrimPrefix(filePath, "./")
		cleanPath = strings.TrimPrefix(cleanPath, "/")

		fullPath := filepath.Join(projectPath, cleanPath)

		// Create directory if needed
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}

		// Write file
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return nil, fmt.Errorf("failed to write file: %w", err)
		}

		createdFiles = append(createdFiles, cleanPath)
		log.Printf("[INFO] Created file: %s", cleanPath)
	}

	log.Printf("[INFO] Successfully generated %d files", len(createdFiles))
	return createdFiles, nil
}

func (g *FileGenerator) extractCodeBlocks(response string) map[string]string {
	files := make(map[string]string)

	lines := strings.Split(response, "\n")
	var currentFile string
	var currentContent []string
	inCodeBlock := false

	for _, line := range lines {
		// Check for file path in markdown code block
		if strings.HasPrefix(line, "```") {
			if !inCodeBlock {
				// Start of code block - try to extract file name
				remaining := strings.TrimPrefix(line, "```")
				remaining = strings.TrimSpace(remaining)

				// Look for filename pattern
				if strings.Contains(remaining, ".") {
					currentFile = strings.TrimSpace(remaining)
					inCodeBlock = true
					currentContent = []string{}
				} else {
					// No filename specified, start code block without filename
					currentFile = fmt.Sprintf("file_%d", len(files)+1)
					inCodeBlock = true
					currentContent = []string{}
				}
			} else {
				// End of code block
				if currentFile != "" && len(currentContent) > 0 {
					files[currentFile] = strings.Join(currentContent, "\n")
					currentFile = ""
				}
				inCodeBlock = false
			}
		} else if inCodeBlock {
			currentContent = append(currentContent, line)
		}
	}

	return files
}

func (g *FileGenerator) CreateProjectDirectory(basePath, projectName string) (string, error) {
	projectPath := filepath.Join(basePath, projectName)

	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create project directory: %w", err)
	}

	return projectPath, nil
}
