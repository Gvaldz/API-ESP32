package server

import (
	temperatureRouters "esp32/src/internal/temperatura/infrastructure"
	motionRouters "esp32/src/internal/motion/infrastructure"
	humidityRouters "esp32/src/internal/humidity/infrastructure"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	temperatureRouters *temperatureRouters.TemperatureRoutes
	motionRouters *motionRouters.MotionRoutes
	humidityRouters *humidityRouters.HumidityRoutes
}

func NewServer(
	tempRoutes *temperatureRouters.TemperatureRoutes,
	motionRoutes *motionRouters.MotionRoutes,
	humidityRoutes *humidityRouters.HumidityRoutes,
) *Server {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	return &Server{
		engine:            r,
		temperatureRouters: tempRoutes,
		motionRouters:      motionRoutes,
		humidityRouters:    humidityRoutes,
	}
}

func (s *Server) Run() error {
	s.temperatureRouters.AttachRoutes(s.engine)
	s.motionRouters.AttachRoutes(s.engine)
	s.humidityRouters.AttachRoutes(s.engine)
	return s.engine.Run(":8080")
}