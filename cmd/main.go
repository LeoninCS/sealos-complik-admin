package main

import (
	"fmt"

	"sealos-complik-admin/internal/config"
	"sealos-complik-admin/internal/database"
	"sealos-complik-admin/internal/logger"
	"sealos-complik-admin/internal/router"
)

const (
	configFile = "configs/config.env.yaml"
)

func main() {
	// Load config

	cfg := config.LoadConfig(configFile)

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

	// Initialize router
	srv := router.InitRouter()
	addr := fmt.Sprintf(":%d", cfg.Port)

	// Start server
	appLogger.Printf("server listening on %s", addr)
	if err := srv.Run(addr); err != nil {
		appLogger.Fatalf("run server: %v", err)
	}
}
