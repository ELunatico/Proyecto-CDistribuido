package routes

import (
	"github.com/gin-gonic/gin"

	"monolito-go/controllers" 
)

func ConfigurarRutas() *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.GET("/", controllers.VerDashboard)

	router.POST("/optimizar", controllers.IniciarOptimizacion)

	return router
}