package main

import (
	"fmt"
	"os"
	"strings"

	"sealos-complik-admin/internal/infra/config"
	"sealos-complik-admin/internal/infra/database"
	"sealos-complik-admin/internal/infra/logger"
	"sealos-complik-admin/internal/infra/migration"
	"sealos-complik-admin/internal/router"
)

const (
	defaultConfigFile  = "/config/config.yaml"
	fallbackConfigFile = "configs/config.yaml"
)

func resolveConfigFile() string {
	if value := strings.TrimSpace(os.Getenv("CONFIG_FILE")); value != "" {
		return value
	}

	if _, err := os.Stat(defaultConfigFile); err == nil {
		return defaultConfigFile
	}

	return fallbackConfigFile
}

func main() {
	// Load config

	cfg := config.LoadConfig(resolveConfigFile())

	// Initialize logger
	appLogger, err := logger.New(cfg.LogDir)
	if err != nil {
		panic(err)
	}

	defer appLogger.CloseWithReport()

	// Initialize database connection
	if _, err := database.Init(cfg.Database); err != nil {
		appLogger.Fatalf("initialize database: %v", err)
	}
	defer database.CloseWithReport(appLogger.Printf)

	if err := migration.AutoMigrate(database.Get()); err != nil {
		appLogger.Fatalf("auto migrate tables: %v", err)
	}

	// Initialize router
	srv := router.InitRouter()
	addr := fmt.Sprintf(":%d", cfg.Port)

	// Start server
	appLogger.Printf("server listening on %s", addr)
	if err := srv.Run(addr); err != nil {
		appLogger.Fatalf("run server: %v", err)
	}
}
