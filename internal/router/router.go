package router

import (
	"sealos-complik-admin/internal/modules/ban"
	"sealos-complik-admin/internal/modules/commitment"
	"sealos-complik-admin/internal/modules/complikviolation"
	"sealos-complik-admin/internal/modules/procscanviolation"
	"sealos-complik-admin/internal/modules/projectconfig"
	"sealos-complik-admin/internal/modules/unban"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	g := gin.Default()
	g.GET("/health", HealthCheck)

	ban.InitBanRoutes(g)
	complikviolation.InitRoutes(g)
	commitment.InitCommitmentRoutes(g)
	projectconfig.InitProjectConfigRoutes(g)
	procscanviolation.InitRoutes(g)
	unban.InitUnbanRoutes(g)

	return g
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "All is well",
	})
}
