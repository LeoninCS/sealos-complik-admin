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
	g.DELETE("/api/violations/:user_id", handler.DeleteViolations)
	g.GET("/api/violations/:user_id", handler.GetViolations)
	g.GET("/api/violations", handler.ListViolations)
	g.GET("/api/users/:user_id/violations-status", handler.GetViolationStatus)
}
