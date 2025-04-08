package infrastructure

import (
	"esp32/src/internal/sensores/humidity/infrastructure/controllers"

	"github.com/gin-gonic/gin"
)

type HumidityRoutes struct {
	CreateHumidityController *controllers.CreateHumidityController
	GetByHamsterController   *controllers.GetByHamsterController
}

func NewHumidityRoutes(
	createHumidityController *controllers.CreateHumidityController,
	getByHamsterController *controllers.GetByHamsterController,
) *HumidityRoutes {
	return &HumidityRoutes{
		CreateHumidityController: createHumidityController,
		GetByHamsterController:   getByHamsterController,
	}
}

func (r *HumidityRoutes) AttachRoutes(router *gin.Engine) {
	humidityGroup := router.Group("/humidity")
	{
		humidityGroup.POST("", r.CreateHumidityController.Create)
		humidityGroup.GET("/hamster/:idHamster", r.GetByHamsterController.GetByHamster)
	}
}
