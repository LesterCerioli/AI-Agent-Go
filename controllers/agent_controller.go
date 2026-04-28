package controllers

import (
	"log"
	"net/http"

	"ai-agent/models"
	"ai-agent/services/implementations"

	"github.com/gofiber/fiber/v2"
)

type AgentController struct {
	agentService *implementations.AgentService
}

func NewAgentController(agentService *implementations.AgentService) *AgentController {
	return &AgentController{
		agentService: agentService,
	}
}

// GenerateCode godoc
// @Summary Generate code using AI
// @Description Generates code based on prompt and creates files in current VS Code workspace
// @Tags Agent
// @Accept json
// @Produce json
// @Param request body models.GenerateCodeRequest true "Code generation request"
// @Success 200 {object} models.GenerateCodeResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /agent/generate [post]
func (c *AgentController) GenerateCode(ctx *fiber.Ctx) error {
	var req models.GenerateCodeRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Printf("[ERROR] Failed to parse request body: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate required fields
	if req.Prompt == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "prompt is required",
		})
	}

	// Set default project name if not provided
	if req.ProjectName == "" {
		req.ProjectName = "generated_project"
	}

	// Call service to generate code
	response, err := c.agentService.GenerateCode(ctx.Context(), &req)
	if err != nil {
		log.Printf("[ERROR] Failed to generate code: %v", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to generate code",
			"details": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

// GetCurrentContext godoc
// @Summary Get current VS Code workspace context
// @Description Returns information about the currently open VS Code workspace
// @Tags Agent
// @Produce json
// @Success 200 {object} models.GetCurrentContextResponse
// @Failure 500 {object} map[string]string
// @Router /agent/context [get]
func (c *AgentController) GetCurrentContext(ctx *fiber.Ctx) error {
	response, err := c.agentService.GetCurrentContext(ctx.Context())
	if err != nil {
		log.Printf("[ERROR] Failed to get current context: %v", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get current context",
			"details": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Returns the health status of the service
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (c *AgentController) HealthCheck(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"service": "ai-agent",
	})
}
