package initializers

import (
	"log"

	"ai-agent/models"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	log.Println("🔄 Starting database migrations...")

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Printf("⚠️ Warning: Error enabling uuid-ossp extension: %v", err)
	}

	if err := db.AutoMigrate(
		&models.Project{},
		&models.Generation{},
	); err != nil {
		log.Printf("⚠️ Warning: Error running migrations: %v", err)
	} else {
		log.Println("✅ Database migrations completed successfully.")
	}

	log.Println("✅ Migrations completed")
}
