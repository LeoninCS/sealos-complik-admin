package commitment

import (
	"sealos-complik-admin/internal/infra/database"

	"github.com/gin-gonic/gin"
)

// InitCommitmentRoutes wires module dependencies and registers commitment APIs.
func InitCommitmentRoutes(g *gin.Engine) {
	repository := NewRepository(database.Get())
	service := NewService(repository)
	handler := NewHandler(service)

	g.POST("/api/commitments", handler.CreateCommitment)
	g.DELETE("/api/commitments/:namespace", handler.DeleteCommitment)
	g.PUT("/api/commitments/:namespace", handler.UpdateCommitment)
	g.GET("/api/commitments/:namespace", handler.GetCommitment)
	g.GET("/api/commitments", handler.ListCommitments)
}
