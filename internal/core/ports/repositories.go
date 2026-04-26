package ports

import (
	"context"

	"github.com/LesterCerioli/AI-Agent-Go/internal/core/domain"
	"github.com/google/uuid"
)

type ProjectRepository interface {
	SaveRequirement(ctx context.Context, requirement *domain.ProjectRequirement) error
	GetRequirementByID(ctx context.Context, id uuid.UUID) (*domain.ProjectRequirement, error)

	SaveArchitecture(ctx context.Context, architecture *domain.ProjectArchitecture) error
	SaveGeneratedFiles(ctx context.Context, files []domain.GeneratedFile) error
}
