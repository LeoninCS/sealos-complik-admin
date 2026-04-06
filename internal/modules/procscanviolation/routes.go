package procscanviolation

import (
	"sealos-complik-admin/internal/infra/database"

	"github.com/gin-gonic/gin"
)

func InitRoutes(g *gin.Engine) {
	repository := NewRepository(database.Get())
	service := NewService(repository)
	handler := NewHandler(service)

	g.POST("/api/procscan-violations", handler.CreateViolation)
	g.PUT("/api/procscan-violations/:id/status", handler.UpdateViolationStatus)
	g.DELETE("/api/procscan-violations/:namespace", handler.DeleteViolations)
	g.GET("/api/procscan-violations/:namespace", handler.GetViolations)
	g.GET("/api/procscan-violations", handler.ListViolations)
	g.GET("/api/namespaces/:namespace/procscan-violations-status", handler.GetViolationStatus)
}
