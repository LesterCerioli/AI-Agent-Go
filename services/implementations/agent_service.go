// services/implementations/agent_service.go
package implementations

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"ai-agent/models"

	"github.com/google/uuid"
)

type AgentService struct {
	db             *sql.DB
	deepSeekClient *DeepSeekClient
	vscodeDetector *VSCodeDetector
	fileGenerator  *FileGenerator
}

func NewAgentService(
	db *sql.DB,
	deepSeekClient *DeepSeekClient,
	vscodeDetector *VSCodeDetector,
	fileGenerator *FileGenerator,
) (*AgentService, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is required")
	}
	if deepSeekClient == nil {
		return nil, fmt.Errorf("deepseek client is required")
	}
	if vscodeDetector == nil {
		return nil, fmt.Errorf("vscode detector is required")
	}
	if fileGenerator == nil {
		return nil, fmt.Errorf("file generator is required")
	}

	return &AgentService{
		db:             db,
		deepSeekClient: deepSeekClient,
		vscodeDetector: vscodeDetector,
		fileGenerator:  fileGenerator,
	}, nil
}

func (s *AgentService) GenerateCode(ctx context.Context, req *models.GenerateCodeRequest) (*models.GenerateCodeResponse, error) {
	log.Printf("[INFO] Starting code generation for prompt: %s", req.Prompt)

	// Detect current workspace
	workspacePath, err := s.vscodeDetector.GetCurrentWorkspace()
	if err != nil {
		return nil, fmt.Errorf("failed to detect workspace: %w", err)
	}

	// Use project name or generate one
	projectName := req.ProjectName
	if projectName == "" {
		projectName = fmt.Sprintf("project_%d", time.Now().Unix())
	}

	// Create project directory
	projectPath, err := s.fileGenerator.CreateProjectDirectory(workspacePath, projectName)
	if err != nil {
		return nil, fmt.Errorf("failed to create project directory: %w", err)
	}

	// Detect language (use provided or detect)
	language := req.Language
	if language == "" {
		language = s.vscodeDetector.DetectProjectLanguage(workspacePath)
		if language == "unknown" {
			language = "go" // default
		}
	}

	// Create project record in database
	projectID, err := s.createProject(ctx, projectName, req.Description, projectPath, language)
	if err != nil {
		log.Printf("[WARN] Failed to create project record: %v", err)
	}

	// Build context for DeepSeek
	contextInfo := fmt.Sprintf("Project: %s\nLanguage: %s\nDescription: %s",
		projectName, language, req.Description)

	// Generate code using DeepSeek
	deepSeekResponse, err := s.deepSeekClient.GenerateCode(ctx, req.Prompt, contextInfo)
	if err != nil {
		s.updateGenerationStatus(ctx, projectID, "failed", err.Error())
		return nil, fmt.Errorf("DeepSeek generation failed: %w", err)
	}

	// Parse and generate files
	filesCreated, err := s.fileGenerator.ParseAndGenerateFiles(projectPath, deepSeekResponse)
	if err != nil {
		s.updateGenerationStatus(ctx, projectID, "failed", err.Error())
		return nil, fmt.Errorf("file generation failed: %w", err)
	}

	// Save generation record
	filesJSON, _ := json.Marshal(filesCreated)
	s.saveGeneration(ctx, projectID, req.Prompt, deepSeekResponse, string(filesJSON), "completed")

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

func (s *AgentService) saveGeneration(ctx context.Context, projectID uuid.UUID, prompt, deepseekResponse, filesCreated, status string) {
	query := `
		INSERT INTO generations (project_id, prompt, deepseek_response, files_created, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`
	s.db.ExecContext(ctx, query, projectID, prompt, deepseekResponse, filesCreated, status)
}

func (s *AgentService) updateGenerationStatus(ctx context.Context, projectID uuid.UUID, status, errorMessage string) {
	query := `
		UPDATE generations 
		SET status = $1, error_message = $2, updated_at = NOW()
		WHERE project_id = $3 AND created_at = (SELECT MAX(created_at) FROM generations WHERE project_id = $3)
	`
	s.db.ExecContext(ctx, query, status, errorMessage, projectID)
}
