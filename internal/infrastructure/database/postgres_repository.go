package database

import (
	"context"

	"github.com/LesterCerioli/AI-Agent-Go/internal/core/domain"
	"github.com/LesterCerioli/AI-Agent-Go/internal/core/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) ports.ProjectRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) SaveRequirement(ctx context.Context, requirement *domain.ProjectRequirement) error {
	return r.db.WithContext(ctx).Save(requirement).Error
}

func (r *postgresRepository) GetRequirementByID(ctx context.Context, id uuid.UUID) (*domain.ProjectRequirement, error) {
	var req domain.ProjectRequirement
	err := r.db.WithContext(ctx).Preload("Architecture.Files").First(&req, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *postgresRepository) SaveArchitecture(ctx context.Context, architecture *domain.ProjectArchitecture) error {
	return r.db.WithContext(ctx).Save(architecture).Error
}

func (r *postgresRepository) SaveGeneratedFiles(ctx context.Context, files []domain.GeneratedFile) error {
	if len(files) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&files).Error
}
