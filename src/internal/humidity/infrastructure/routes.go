package infrastructure

import (
	"github.com/gin-gonic/gin"
	"esp32/src/internal/humidity/infrastructure/controllers"
)

type HumidityRoutes struct {
	CreateHumidityController *controllers.CreateHumidityController
	GetByHamsterController      *controllers.GetByHamsterController
}

func NewHumidityRoutes(
	createHumidityController *controllers.CreateHumidityController,
	getByHamsterController *controllers.GetByHamsterController,
) *HumidityRoutes {
	return &HumidityRoutes{
		CreateHumidityController: createHumidityController,
		GetByHamsterController:      getByHamsterController,
	}
}

func (r *HumidityRoutes) AttachRoutes(router *gin.Engine) {
	humidityGroup := router.Group("/humidity")
	{
		humidityGroup.POST("", r.CreateHumidityController.Create)
		humidityGroup.GET("/hamster/:idHamster", r.GetByHamsterController.GetByHamster)
	}
}
