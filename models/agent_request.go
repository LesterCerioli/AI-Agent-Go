package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Path        string    `gorm:"uniqueIndex" json:"path"`
	Language    string    `json:"language"`
	Status      string    `gorm:"default:'pending'" json:"status"` // pending, generating, completed, failed
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Generation struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ProjectID        uuid.UUID `gorm:"type:uuid;not null" json:"project_id"`
	Prompt           string    `gorm:"type:text;not null" json:"prompt"`
	DeepSeekResponse string    `gorm:"type:text" json:"deepseek_response"`
	FilesCreated     string    `gorm:"type:text" json:"files_created"`  // JSON array
	Status           string    `gorm:"default:'pending'" json:"status"` // pending, completed, failed
	ErrorMessage     string    `gorm:"type:text" json:"error_message"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type GenerateCodeRequest struct {
	Prompt      string `json:"prompt"`
	ProjectName string `json:"project_name"`
	Description string `json:"description"`
	Language    string `json:"language"`
}

type GenerateCodeResponse struct {
	Success      bool     `json:"success"`
	Message      string   `json:"message"`
	ProjectID    string   `json:"project_id"`
	ProjectPath  string   `json:"project_path"`
	FilesCreated []string `json:"files_created"`
}

type GetCurrentContextResponse struct {
	WorkspacePath string   `json:"workspace_path"`
	ProjectName   string   `json:"project_name"`
	Files         []string `json:"files"`
	Language      string   `json:"language"`
}
