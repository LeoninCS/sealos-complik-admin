package router

import (
	"sealos-complik-admin/internal/modules/projectconfig"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	g := gin.Default()
	g.GET("/health", HealthCheck)

	projectconfig.InitProjectConfigRoutes(g)

	return g
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "All is well",
	})
}
