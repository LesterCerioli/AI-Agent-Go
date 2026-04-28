package initializers

import (
	"context"
	"log"
	"os"

	"ai-agent/services/implementations"

	"gorm.io/gorm"
)

type Services struct {
	AgentService   *implementations.AgentService
	VSCodeDetector *implementations.VSCodeDetector
	FileGenerator  *implementations.FileGenerator
	OllamaClient   *implementations.OllamaClient
}

func InitServices(db *gorm.DB) *Services {
	log.Println("Initializing services...")

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB from gorm: %v", err)
	}

	
	ollamaClient, err := implementations.NewOllamaClient()
	if err != nil {
		log.Fatalf("Failed to initialize OllamaClient: %v", err)
	}

	ctx := context.Background()
	if err := ollamaClient.CheckHealth(ctx); err != nil {
		log.Printf("⚠️  Ollama health check failed: %v", err)
		log.Printf("   Make sure Ollama is running: docker-compose up -d")
		log.Printf("   And model is pulled: docker exec ollama ollama pull %s", os.Getenv("OLLAMA_MODEL"))
		log.Fatal("❌ Ollama is not healthy. Cannot continue without Ollama.")
	} else {
		log.Println("✅ OllamaClient initialized and healthy")
		if models, err := ollamaClient.ListModels(ctx); err == nil {
			log.Printf("📦 Available Ollama models: %v", models)
		}
	}

	vscodeDetector, err := implementations.NewVSCodeDetector()
	if err != nil {
		log.Fatalf("Failed to initialize VSCodeDetector: %v", err)
	}
	log.Println("✅ VSCodeDetector initialized")

	fileGenerator, err := implementations.NewFileGenerator()
	if err != nil {
		log.Fatalf("Failed to initialize FileGenerator: %v", err)
	}
	log.Println("✅ FileGenerator initialized")
	
	agentService, err := implementations.NewAgentService(
		sqlDB,
		ollamaClient, // Adicionado OllamaClient como segundo argumento
		vscodeDetector,
		fileGenerator,
	)
	if err != nil {
		log.Fatalf("Failed to initialize AgentService: %v", err)
	}
	log.Println("✅ AgentService initialized")

	log.Println("All services initialized successfully!")

	return &Services{
		AgentService:   agentService,
		VSCodeDetector: vscodeDetector,
		FileGenerator:  fileGenerator,
		OllamaClient:   ollamaClient,
	}
}
