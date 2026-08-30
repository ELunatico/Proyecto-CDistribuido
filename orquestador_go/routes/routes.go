package routes

import (
	"github.com/gin-gonic/gin"
	
	"monolito-go/controllers"
)

func ConfigurarRutas() *gin.Engine {
	router := gin.Default()

	router.LoadHTMLGlob("views/*")

	router.GET("/", controllers.VerDashboard)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"estado": "Orquestador activo",
		})
	})

	router.POST("/optimizar", controllers.IniciarOptimizacion)

	return router
}