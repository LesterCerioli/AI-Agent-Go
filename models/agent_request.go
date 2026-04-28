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
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ProjectID      uuid.UUID `gorm:"type:uuid;not null" json:"project_id"`
	Prompt         string    `gorm:"type:text;not null" json:"prompt"`
	Requirements   string    `gorm:"type:text" json:"requirements"`    // Adicionado: requisitos do projeto
	OllamaResponse string    `gorm:"type:text" json:"ollama_response"` // Mudado de DeepSeekResponse
	FilesCreated   string    `gorm:"type:text" json:"files_created"`   // JSON array
	Status         string    `gorm:"default:'pending'" json:"status"`  // pending, completed, failed
	ErrorMessage   string    `gorm:"type:text" json:"error_message"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type GenerateCodeRequest struct {
	Prompt       string `json:"prompt" validate:"required"`
	ProjectName  string `json:"projectName"`                      // Nome do projeto
	Requirements string `json:"requirements" validate:"required"` // Requisitos técnicos detalhados
	ProjectPath  string `json:"project_path" validate:"required"` // Caminho absoluto onde gerar o código
	Language     string `json:"language"`                         // Opcional: go, python, javascript, etc
	Description  string `json:"description"`                      // Opcional: descrição adicional
}

type GenerateCodeResponse struct {
	Success      bool     `json:"success"`
	Message      string   `json:"message"`
	ProjectID    string   `json:"projectId"`    // Mantido camelCase para consistência
	ProjectPath  string   `json:"projectPath"`  // Mantido camelCase
	FilesCreated []string `json:"filesCreated"` // Mantido camelCase
}

type GetCurrentContextResponse struct {
	WorkspacePath string   `json:"workspace_path"`
	ProjectName   string   `json:"project_name"`
	Files         []string `json:"files"`
	Language      string   `json:"language"`
}
