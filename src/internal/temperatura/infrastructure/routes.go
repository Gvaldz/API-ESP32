package infrastructure

import (
	"github.com/gin-gonic/gin"
	"esp32/src/internal/temperatura/infrastructure/controllers"
)

type TemperatureRoutes struct {
	CreateTemperatureController *controllers.CreateTemperatureController
	GetByHamsterController      *controllers.GetByHamsterController
}

func NewTemperatureRoutes(
	createTemperatureController *controllers.CreateTemperatureController,
	getByHamsterController *controllers.GetByHamsterController,
) *TemperatureRoutes {
	return &TemperatureRoutes{
		CreateTemperatureController: createTemperatureController,
		GetByHamsterController:      getByHamsterController,
	}
}

func (r *TemperatureRoutes) AttachRoutes(router *gin.Engine) {
	temperatureGroup := router.Group("/temperatures")
	{
		temperatureGroup.POST("", r.CreateTemperatureController.Create)
		temperatureGroup.GET("/hamster/:idHamster", r.GetByHamsterController.GetByHamster)
	}
}
