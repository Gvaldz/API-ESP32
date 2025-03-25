package server

import (
	temperatureRouters 	"esp32/src/internal/temperatura/infrastructure"
	// motionRouters 		"esp32/src/internal/motion/infrastructure"
	humidityRouters 	"esp32/src/internal/humidity/infrastructure"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run(
	temperatureRouters *temperatureRouters.TemperatureRoutes,
	// motionRouters *motionRouters.MotionRoutes,
	humidityRouters *humidityRouters.HumidityRoutes,
) {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	
	}))

	temperatureRouters.AttachRoutes(r)
	// motionRouters.AttachRoutes(r)
	humidityRouters.AttachRoutes(r)

	r.Run(":8080")
}
