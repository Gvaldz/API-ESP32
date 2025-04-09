package server

import (
	cagesRouters "esp32/src/internal/sensores/cages/infrastructure"
	foodRouters "esp32/src/internal/sensores/food/infrastructure"
	humidityRouters "esp32/src/internal/sensores/humidity/infrastructure"
	motionRouters "esp32/src/internal/sensores/motion/infrastructure"
	temperatureRouters "esp32/src/internal/sensores/temperatura/infrastructure"
	websocketRouters "esp32/src/internal/services/websocket/infrastructure"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
    engine             *gin.Engine
    temperatureRouters *temperatureRouters.TemperatureRoutes
    motionRouters      *motionRouters.MotionRoutes
    humidityRouters    *humidityRouters.HumidityRoutes
    foodRouters        *foodRouters.FoodRoutes
    cagesRouters       *cagesRouters.CageRoutes
    websocketRouters   *websocketRouters.WebSocketRoutes
}

func NewServer(
    tempRoutes    *temperatureRouters.TemperatureRoutes,
    motionRoutes  *motionRouters.MotionRoutes,
    humidityRoutes *humidityRouters.HumidityRoutes,
    foodRoutes    *foodRouters.FoodRoutes,
    cageRoutes    *cagesRouters.CageRoutes,
    wsRoutes      *websocketRouters.WebSocketRoutes,
)*Server {
    r := gin.Default()

    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"}, 
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
    }))

    return &Server{
        engine:            r,
        temperatureRouters: tempRoutes,
        motionRouters:      motionRoutes,
        humidityRouters:    humidityRoutes,
        foodRouters:        foodRoutes,
        cagesRouters:       cageRoutes,
        websocketRouters:   wsRoutes,
    }
}

func (s *Server) Run() error {
    s.temperatureRouters.AttachRoutes(s.engine)
    s.motionRouters.AttachRoutes(s.engine)
    s.humidityRouters.AttachRoutes(s.engine)
    s.foodRouters.AttachRoutes(s.engine)
    s.cagesRouters.AttachRoutes(s.engine)
    s.websocketRouters.AttachRoutes(s.engine) 
    return s.engine.Run(":8080")
}
