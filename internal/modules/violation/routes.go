package violation

import (
	"sealos-complik-admin/internal/infra/database"

	"github.com/gin-gonic/gin"
)

// InitViolationRoutes wires module dependencies and registers violation APIs.
func InitViolationRoutes(g *gin.Engine) {
	repository := NewRepository(database.Get())
	service := NewService(repository)
	handler := NewHandler(service)

	g.POST("/api/violations", handler.CreateViolation)
	g.DELETE("/api/violations/:namespace", handler.DeleteViolations)
	g.GET("/api/violations/:namespace", handler.GetViolations)
	g.GET("/api/violations", handler.ListViolations)
	g.GET("/api/namespaces/:namespace/violations-status", handler.GetViolationStatus)
}
