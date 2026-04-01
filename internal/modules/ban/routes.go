package ban

import (
	"sealos-complik-admin/internal/infra/database"

	"github.com/gin-gonic/gin"
)

// InitBanRoutes wires module dependencies and registers ban APIs.
func InitBanRoutes(g *gin.Engine) {
	repository := NewRepository(database.Get())
	service := NewService(repository)
	handler := NewHandler(service)

	g.POST("/api/bans", handler.CreateBan)
	g.DELETE("/api/bans/:namespace", handler.DeleteBans)
	g.GET("/api/bans/:namespace", handler.GetBans)
	g.GET("/api/bans", handler.ListBans)
	g.GET("/api/namespaces/:namespace/ban-status", handler.GetBanStatus)
}
