package services

import (
	"context"
	"log"

	"github.com/LesterCerioli/AI-Agent-Go/internal/core/domain"
	"github.com/LesterCerioli/AI-Agent-Go/internal/core/ports"
)

type generationService struct {
	aiClient ports.AIClient
	repo     ports.ProjectRepository
}

func NewGenerationService(aiClient ports.AIClient, repo ports.ProjectRepository) ports.GenerationService {
	return &generationService{
		aiClient: aiClient,
		repo:     repo,
	}
}

func (s *generationService) ProcessRequirement(ctx context.Context, requirement *domain.ProjectRequirement) error {

	requirement.Status = "processing"
	if err := s.repo.SaveRequirement(ctx, requirement); err != nil {
		return err
	}

	go func(req domain.ProjectRequirement) {
		bgCtx := context.Background() // Use background context since the request context might expire

		log.Printf("Starting generation for requirement ID: %s\n", req.ID)

		arch, err := s.aiClient.AnalyzeAndArchitect(bgCtx, req.Prompt)
		if err != nil {
			log.Printf("Error analyzing architecture: %v\n", err)
			s.updateStatus(bgCtx, &req, "failed")
			return
		}

		arch.RequirementID = req.ID
		if err := s.repo.SaveArchitecture(bgCtx, arch); err != nil {
			log.Printf("Error saving architecture: %v\n", err)
			s.updateStatus(bgCtx, &req, "failed")
			return
		}

		files, err := s.aiClient.GenerateCode(bgCtx, arch)
		if err != nil {
			log.Printf("Error generating code: %v\n", err)
			s.updateStatus(bgCtx, &req, "failed")
			return
		}

		for i := range files {
			files[i].ArchitectureID = arch.ID
		}

		if err := s.repo.SaveGeneratedFiles(bgCtx, files); err != nil {
			log.Printf("Error saving generated files: %v\n", err)
			s.updateStatus(bgCtx, &req, "failed")
			return
		}

		s.updateStatus(bgCtx, &req, "completed")
		log.Printf("Completed generation for requirement ID: %s\n", req.ID)
	}(*requirement)

	return nil
}

func (s *generationService) updateStatus(ctx context.Context, req *domain.ProjectRequirement, status string) {
	req.Status = status
	if err := s.repo.SaveRequirement(ctx, req); err != nil {
		log.Printf("Failed to update status to %s for requirement ID %s: %v\n", status, req.ID, err)
	}
}
