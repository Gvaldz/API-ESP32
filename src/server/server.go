package server

import (
	temperatureRouters 	"esp32/src/internal/temperatura/infrastructure"
	motionRouters 		"esp32/src/internal/motion/infrastructure"
	humidityRouters 	"esp32/src/internal/humidity/infrastructure"
	foodRouters 		"esp32/src/internal/food/infrastructure"
	usersRouters		"esp32/src/internal/users/infrastructure"
	cagesRouters		"esp32/src/internal/cages/infrastructure"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine 				 *gin.Engine
	temperatureRouters 	 *temperatureRouters.TemperatureRoutes
	motionRouters 		 *motionRouters.MotionRoutes
	humidityRouters 	 *humidityRouters.HumidityRoutes
	foodRouters 		 *foodRouters.FoodRoutes
	usersRouters 		 *usersRouters.UserRoutes
	cagesRouters 		 *cagesRouters.CageRoutes
}

func NewServer(
	tempRoutes 			 *temperatureRouters.TemperatureRoutes,
	motionRoutes         *motionRouters.MotionRoutes,
	humidityRoutes 		 *humidityRouters.HumidityRoutes,
	foodRoutes 	  		 *foodRouters.FoodRoutes,
	userRoutes			 *usersRouters.UserRoutes,
	cageRoutes 			 *cagesRouters.CageRoutes, 
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
		foodRouters: 	  	foodRoutes,
		usersRouters: 		userRoutes,
		cagesRouters: 		cageRoutes,	
	}
}

func (s *Server) Run() error {
	s.temperatureRouters.AttachRoutes(s.engine)
	s.motionRouters.AttachRoutes(s.engine)
	s.humidityRouters.AttachRoutes(s.engine)
	s.foodRouters.AttachRoutes(s.engine)
	s.usersRouters.AttachRoutes(s.engine)
	s.cagesRouters.AttachRoutes(s.engine)
	return s.engine.Run(":8080")
}