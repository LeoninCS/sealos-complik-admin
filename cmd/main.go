package main

import (
	"fmt"

	"sealos-complik-admin/internal/infra/config"
	"sealos-complik-admin/internal/infra/database"
	"sealos-complik-admin/internal/infra/logger"
	"sealos-complik-admin/internal/modules/projectconfig"
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

	if err := projectconfig.AutoMigrate(database.Get()); err != nil {
		appLogger.Fatalf("auto migrate project config table: %v", err)
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
