package implementations

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"ai-agent/models"

	"github.com/google/uuid"
)

type AgentService struct {
	db             *sql.DB
	ollamaClient   *OllamaClient
	vscodeDetector *VSCodeDetector
	fileGenerator  *FileGenerator
}

func NewAgentService(
	db *sql.DB,
	ollamaClient *OllamaClient,
	vscodeDetector *VSCodeDetector,
	fileGenerator *FileGenerator,
) (*AgentService, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is required")
	}
	if ollamaClient == nil {
		return nil, fmt.Errorf("ollama client is required")
	}
	if vscodeDetector == nil {
		return nil, fmt.Errorf("vscode detector is required")
	}
	if fileGenerator == nil {
		return nil, fmt.Errorf("file generator is required")
	}

	return &AgentService{
		db:             db,
		ollamaClient:   ollamaClient,
		vscodeDetector: vscodeDetector,
		fileGenerator:  fileGenerator,
	}, nil
}

func (s *AgentService) GenerateCode(ctx context.Context, req *models.GenerateCodeRequest) (*models.GenerateCodeResponse, error) {
	log.Printf("[INFO] Starting code generation for prompt: %s", req.Prompt)
	log.Printf("[INFO] Project path: %s", req.ProjectPath)
	log.Printf("[INFO] Requirements: %s", req.Requirements)

	if req.ProjectPath == "" {
		return nil, fmt.Errorf("project_path is required")
	}

	projectPath := req.ProjectPath

	if err := s.fileGenerator.EnsureProjectDirectory(projectPath); err != nil {
		return nil, fmt.Errorf("failed to create or access project directory: %w", err)
	}

	projectName := req.ProjectName
	if projectName == "" {
		projectName = filepath.Base(projectPath)
	}

	language := req.Language
	if language == "" {
		language = s.vscodeDetector.DetectProjectLanguage(projectPath)
		if language == "unknown" {
			language = "go" // default
		}
	}

	projectID, err := s.createProject(ctx, projectName, req.Description, projectPath, language)
	if err != nil {
		log.Printf("[WARN] Failed to create project record: %v", err)
	}

	contextInfo := fmt.Sprintf(`Project Name: %s
Language: %s
User Prompt: %s
Technical Requirements: %s

Important: Generate production-ready code following best practices for %s.
Return ONLY valid JSON with file paths as keys and file contents as values.`,
		projectName, language, req.Prompt, req.Requirements, language)

	ollamaResponse, err := s.ollamaClient.GenerateCode(ctx, req.Prompt, contextInfo)
	if err != nil {
		s.updateGenerationStatus(ctx, projectID, "failed", err.Error())
		return nil, fmt.Errorf("Ollama generation failed: %w", err)
	}

	filesCreated, err := s.fileGenerator.ParseAndGenerateFiles(projectPath, ollamaResponse)
	if err != nil {
		s.updateGenerationStatus(ctx, projectID, "failed", err.Error())
		return nil, fmt.Errorf("file generation failed: %w", err)
	}

	filesJSON, _ := json.Marshal(filesCreated)
	s.saveGeneration(ctx, projectID, req.Prompt, req.Requirements, ollamaResponse, string(filesJSON), "completed")

	log.Printf("[INFO] Code generation completed. Created %d files in %s", len(filesCreated), projectPath)

	return &models.GenerateCodeResponse{
		Success:      true,
		Message:      fmt.Sprintf("Successfully generated %d files", len(filesCreated)),
		ProjectID:    projectID.String(),
		ProjectPath:  projectPath,
		FilesCreated: filesCreated,
	}, nil
}

func (s *AgentService) GetCurrentContext(ctx context.Context) (*models.GetCurrentContextResponse, error) {
	workspacePath, err := s.vscodeDetector.GetCurrentWorkspace()
	if err != nil {
		return nil, fmt.Errorf("failed to detect workspace: %w", err)
	}

	files, err := s.vscodeDetector.ScanProjectFiles(workspacePath)
	if err != nil {
		files = []string{}
	}

	language := s.vscodeDetector.DetectProjectLanguage(workspacePath)

	return &models.GetCurrentContextResponse{
		WorkspacePath: workspacePath,
		ProjectName:   filepath.Base(workspacePath),
		Files:         files,
		Language:      language,
	}, nil
}

func (s *AgentService) createProject(ctx context.Context, name, description, path, language string) (uuid.UUID, error) {
	query := `
		INSERT INTO projects (name, description, path, language, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id
	`

	var id uuid.UUID
	err := s.db.QueryRowContext(ctx, query, name, description, path, language, "active").Scan(&id)
	return id, err
}

func (s *AgentService) saveGeneration(ctx context.Context, projectID uuid.UUID, prompt, requirements, ollamaResponse, filesCreated, status string) {
	query := `
		INSERT INTO generations (project_id, prompt, requirements, ollama_response, files_created, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`
	s.db.ExecContext(ctx, query, projectID, prompt, requirements, ollamaResponse, filesCreated, status)
}

func (s *AgentService) updateGenerationStatus(ctx context.Context, projectID uuid.UUID, status, errorMessage string) {
	if projectID == uuid.Nil {
		return
	}
	query := `
		UPDATE generations 
		SET status = $1, error_message = $2, updated_at = NOW()
		WHERE project_id = $3 AND created_at = (SELECT MAX(created_at) FROM generations WHERE project_id = $3)
	`
	s.db.ExecContext(ctx, query, status, errorMessage, projectID)
}
