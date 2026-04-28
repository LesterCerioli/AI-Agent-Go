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

func (g *FileGenerator) EnsureProjectDirectory(projectPath string) error {
	log.Printf("[INFO] Ensuring project directory exists: %s", projectPath)

	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		if err := os.MkdirAll(projectPath, 0755); err != nil {
			return fmt.Errorf("failed to create project directory: %w", err)
		}
		log.Printf("[INFO] Created project directory: %s", projectPath)
	} else if err != nil {
		return fmt.Errorf("failed to check project directory: %w", err)
	} else {
		log.Printf("[INFO] Project directory already exists: %s", projectPath)
	}

	return nil
}

func (g *FileGenerator) ParseAndGenerateFiles(projectPath string, ollamaResponse string) ([]string, error) {
	log.Printf("[INFO] Parsing Ollama response and generating files in: %s", projectPath)

	var files map[string]string

	if err := json.Unmarshal([]byte(ollamaResponse), &files); err != nil {

		files = g.extractCodeBlocks(ollamaResponse)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files found in Ollama response")
	}

	var createdFiles []string

	for filePath, content := range files {
		// Limpa o path
		cleanPath := strings.TrimPrefix(filePath, "./")
		cleanPath = strings.TrimPrefix(cleanPath, "/")
		cleanPath = filepath.Clean(cleanPath)

		fullPath := filepath.Join(projectPath, cleanPath)

		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return nil, fmt.Errorf("failed to write file %s: %w", cleanPath, err)
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

				if remaining != "" && (strings.Contains(remaining, ".") || strings.Contains(remaining, "/")) {
					currentFile = remaining
				} else {
					currentFile = fmt.Sprintf("file_%d", len(files)+1)
				}
				inCodeBlock = true
				currentContent = []string{}
			} else {
				// Fim do code block
				if currentFile != "" && len(currentContent) > 0 {
					files[currentFile] = strings.Join(currentContent, "\n")
				}
				currentFile = ""
				inCodeBlock = false
			}
		} else if inCodeBlock {
			currentContent = append(currentContent, line)
		}
	}

	return files
}
