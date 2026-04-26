package controllers

import (
	"github.com/LesterCerioli/AI-Agent-Go/internal/core/domain"
	"github.com/LesterCerioli/AI-Agent-Go/internal/core/ports"
	"github.com/gofiber/fiber/v2"
)

type GenerationController struct {
	genService ports.GenerationService
}

func NewGenerationController(genService ports.GenerationService) *GenerationController {
	return &GenerationController{
		genService: genService,
	}
}

type GenerateRequest struct {
	Prompt string `json:"prompt"`
}

func (c *GenerationController) GenerateProject(ctx *fiber.Ctx) error {
	var req GenerateRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Prompt == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Prompt is required",
		})
	}

	requirement := &domain.ProjectRequirement{
		Prompt: req.Prompt,
		Status: "pending",
	}

	err := c.genService.ProcessRequirement(ctx.Context(), requirement)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process requirement",
		})
	}

	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message":        "Generation started",
		"requirement_id": requirement.ID,
	})
}
