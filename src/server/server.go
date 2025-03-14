package server

import (
	temperatureRouters 	"esp32/src/internal/temperatura/infrastructure"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run(
	temperatureRouters *temperatureRouters.TemperatureRoutes,
) {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	
	}))

	temperatureRouters.AttachRoutes(r)

	r.Run(":8080")
}
