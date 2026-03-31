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
	g.DELETE("/api/commitments/:user_id", handler.DeleteCommitment)
	g.PUT("/api/commitments/:user_id", handler.UpdateCommitment)
	g.GET("/api/commitments/:user_id", handler.GetCommitment)
	g.GET("/api/commitments", handler.ListCommitments)
}
