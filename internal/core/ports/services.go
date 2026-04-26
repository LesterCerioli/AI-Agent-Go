package ports

import (
	"context"

	"github.com/LesterCerioli/AI-Agent-Go/internal/core/domain"
)

type AIClient interface {
	AnalyzeAndArchitect(ctx context.Context, prompt string) (*domain.ProjectArchitecture, error)
	GenerateCode(ctx context.Context, architecture *domain.ProjectArchitecture) ([]domain.GeneratedFile, error)
}

type GenerationService interface {
	ProcessRequirement(ctx context.Context, requirement *domain.ProjectRequirement) error
}
