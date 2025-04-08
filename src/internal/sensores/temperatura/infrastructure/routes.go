package infrastructure

import (
	"esp32/src/internal/sensores/temperatura/infrastructure/controllers"

	"github.com/gin-gonic/gin"
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
