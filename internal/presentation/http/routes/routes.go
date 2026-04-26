package routes

import (
	"github.com/LesterCerioli/AI-Agent-Go/internal/presentation/http/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, genController *controllers.GenerationController) {
	api := app.Group("/api/v1")

	api.Post("/generate", genController.GenerateProject)
}
