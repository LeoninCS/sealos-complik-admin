package router

import "github.com/gin-gonic/gin"

func InitRouter() *gin.Engine {
	g := gin.Default()
	g.GET("/health", HealthCheck)
	return g
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "All is well",
	})
}
