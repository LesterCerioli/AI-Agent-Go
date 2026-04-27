package initializers

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitialDB() *gorm.DB {
	dbHost := getEnv("DB_HOST", "")
	dbUser := getEnv("DB_USER", "")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "")
	dbPort := getEnv("DB_PORT", "5432")
	dbTimezone := getEnv("DB_TIMEZONE", "UTC")
	sslMode := getEnv("DB_SSL_MODE", "require")

	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatal("Missing required database configuration in environment variables")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s default_query_exec_mode=simple_protocol",
		dbHost, dbUser, dbPassword, dbName, dbPort, sslMode, dbTimezone,
	)

	maskedDsn := maskSensitiveInfo(dsn, dbPassword)
	log.Printf("Connecting to production database with DSN: %s", maskedDsn)

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Warn,
			Colorful:      false,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:               gormLogger,
		DisableAutomaticPing: false,
		PrepareStmt:          false, // redundante, mas mantemos
	})
	if err != nil {
		log.Fatalf("Error connecting to the production database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error getting DB instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(0)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(30 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatalf("Error pinging the production database: %v", err)
	}

	log.Println("Successfully connected to production database")
	return db
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func maskSensitiveInfo(dsn, secret string) string {
	if secret != "" {
		return strings.Replace(dsn, secret, "*****", -1)
	}
	return dsn
}
