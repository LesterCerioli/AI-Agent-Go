package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/LesterCerioli/AI-Agent-Go/internal/application/services"
	"github.com/LesterCerioli/AI-Agent-Go/internal/core/domain"
	"github.com/LesterCerioli/AI-Agent-Go/internal/infrastructure/ai"
	"github.com/LesterCerioli/AI-Agent-Go/internal/infrastructure/database"
	"github.com/LesterCerioli/AI-Agent-Go/internal/presentation/http/controllers"
	"github.com/LesterCerioli/AI-Agent-Go/internal/presentation/http/routes"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	dsn := "host=localhost user=postgres password=postgres dbname=ai_agent port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database. Starting without DB for demo purposes. Error: %v\n", err)
	} else {

		err = db.AutoMigrate(
			&domain.ProjectRequirement{},
			&domain.ProjectArchitecture{},
			&domain.GeneratedFile{},
		)
		if err != nil {
			log.Fatalf("Failed to auto migrate database schemas: %v", err)
		}
	}

	repo := database.NewPostgresRepository(db)
	aiClient := ai.NewMockAIClient()

	genService := services.NewGenerationService(aiClient, repo)

	genController := controllers.NewGenerationController(genService)

	app := fiber.New(fiber.Config{
		AppName: "AI-Driven Code Generation Engine v1.0",
	})

	app.Use(logger.New())
	app.Use(recover.New())

	routes.SetupRoutes(app, genController)

	log.Println("Starting AI-Driven Code Generation Engine on port 3000...")
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
