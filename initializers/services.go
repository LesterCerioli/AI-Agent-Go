// initializers/services.go
package initializers

import (
	"log"
	"os"

	"ai-agent/services/implementations"

	"gorm.io/gorm"
)

type Services struct {
	AgentService   *implementations.AgentService
	DeepSeekClient *implementations.DeepSeekClient
	VSCodeDetector *implementations.VSCodeDetector
	FileGenerator  *implementations.FileGenerator
}

func InitServices(db *gorm.DB) *Services {
	log.Println("Initializing services...")

	// Get sql.DB from gorm
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB from gorm: %v", err)
	}

	// 1. Initialize DeepSeek Client
	deepSeekKey := os.Getenv("DEEPSEEK_KEY")
	if deepSeekKey == "" {
		log.Fatal("DEEPSEEK_KEY environment variable is not set")
	}

	deepSeekURL := os.Getenv("DEEPSEEK_URL")
	if deepSeekURL == "" {
		deepSeekURL = "https://api.deepseek.com/v1"
	}

	deepSeekClient, err := implementations.NewDeepSeekClient(deepSeekKey, deepSeekURL)
	if err != nil {
		log.Fatalf("Failed to initialize DeepSeekClient: %v", err)
	}
	log.Println("✅ DeepSeekClient initialized")

	// 2. Initialize VSCode Detector
	vscodeDetector, err := implementations.NewVSCodeDetector()
	if err != nil {
		log.Fatalf("Failed to initialize VSCodeDetector: %v", err)
	}
	log.Println("✅ VSCodeDetector initialized")

	// 3. Initialize File Generator
	fileGenerator, err := implementations.NewFileGenerator()
	if err != nil {
		log.Fatalf("Failed to initialize FileGenerator: %v", err)
	}
	log.Println("✅ FileGenerator initialized")

	// 4. Initialize Agent Service (with all dependencies)
	agentService, err := implementations.NewAgentService(
		sqlDB,
		deepSeekClient,
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
		DeepSeekClient: deepSeekClient,
		VSCodeDetector: vscodeDetector,
		FileGenerator:  fileGenerator,
	}
}
