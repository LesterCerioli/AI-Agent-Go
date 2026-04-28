package main

import (
	"log"
	"os"

	"ai-agent/controllers"
	"ai-agent/initializers"

	_ "ai-agent/cmd/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

func initializeRoutes(app *fiber.App, services *initializers.Services) {

	agentController := controllers.NewAgentController(services.AgentService)

	api := app.Group("/api/v1")

	// Agent routes (public - no auth required for now)
	agent := api.Group("/agent")
	agent.Post("/generate", agentController.GenerateCode)
	agent.Get("/context", agentController.GetCurrentContext)

	// Health check
	api.Get("/health", agentController.HealthCheck)
}

func configureMiddleware(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path}\n",
	}))
}

// @title AI Agent API
// @version 1.0
// @description AI Agent with DeepSeek integration for automatic code generation in VS Code
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@aiagent.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api/v1
// @schemes http https

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Initialize database
	db := initializers.InitialDB()

	// Run migrations
	initializers.RunMigrations(db)

	// Initialize services
	services := initializers.InitServices(db)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "AI Agent API v1",
	})

	// Configure middleware
	configureMiddleware(app)

	// Swagger documentation setup
	swaggerUser := os.Getenv("SWAGGER_USER")
	swaggerPass := os.Getenv("SWAGGER_PASSWORD")

	// Swagger configuration with optional authentication
	if swaggerUser == "" || swaggerPass == "" {
		log.Println("⚠️  Warning: SWAGGER_USER and SWAGGER_PASSWORD not set, Swagger will be public")
		app.Get("/swagger/*", swagger.HandlerDefault)
		app.Get("/docs", func(c *fiber.Ctx) error {
			return c.Redirect("/swagger/index.html")
		})
	} else {
		// Swagger with authentication
		swaggerAuth := basicauth.New(basicauth.Config{
			Users: map[string]string{
				swaggerUser: swaggerPass,
			},
			Realm: "Swagger Restricted",
		})

		app.Get("/swagger/*", swaggerAuth, swagger.HandlerDefault)
		app.Get("/docs", swaggerAuth, func(c *fiber.Ctx) error {
			return c.Redirect("/swagger/index.html")
		})
	}

	// Initialize routes
	initializeRoutes(app, services)

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	// Start server
	log.Printf("🚀 AI Agent Server is running on http://0.0.0.0:%s", port)
	log.Printf("📚 Swagger UI available at: http://localhost:%s/swagger/index.html", port)
	log.Printf("🔧 API endpoints available at: http://localhost:%s/api/v1", port)
	log.Printf("🤖 Agent endpoints:")
	log.Printf("   POST   /api/v1/agent/generate - Generate code with AI")
	log.Printf("   GET    /api/v1/agent/context  - Get current VS Code workspace context")
	log.Printf("   GET    /api/v1/health         - Health check")

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
