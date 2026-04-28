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

	if err := json.Unmarshal([]byte(deepseekResponse), &files); err != nil {

		files = g.extractCodeBlocks(deepseekResponse)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files found in DeepSeek response")
	}

	var createdFiles []string

	for filePath, content := range files {

		cleanPath := strings.TrimPrefix(filePath, "./")
		cleanPath = strings.TrimPrefix(cleanPath, "/")

		fullPath := filepath.Join(projectPath, cleanPath)

		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}

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

		if strings.HasPrefix(line, "```") {
			if !inCodeBlock {

				remaining := strings.TrimPrefix(line, "```")
				remaining = strings.TrimSpace(remaining)

				if strings.Contains(remaining, ".") {
					currentFile = strings.TrimSpace(remaining)
					inCodeBlock = true
					currentContent = []string{}
				} else {

					currentFile = fmt.Sprintf("file_%d", len(files)+1)
					inCodeBlock = true
					currentContent = []string{}
				}
			} else {

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
