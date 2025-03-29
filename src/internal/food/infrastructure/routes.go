package infrastructure

import (
	"github.com/gin-gonic/gin"
	"esp32/src/internal/food/infrastructure/controllers"
)

type FoodRoutes struct {
	CreateStatusFoodController *controllers.CreateStatusFoodController
	GetByHamsterController      *controllers.GetByHamsterController
}

func NewFoodRoutes(
	createStatusFoodController *controllers.CreateStatusFoodController,
	getByHamsterController *controllers.GetByHamsterController,
) *FoodRoutes {
	return &FoodRoutes{
		CreateStatusFoodController: createStatusFoodController,
		GetByHamsterController:      getByHamsterController,
	}
}

func (r *FoodRoutes) AttachRoutes(router *gin.Engine) {
	temperatureGroup := router.Group("/food")
	{
		temperatureGroup.POST("", r.CreateStatusFoodController.Create)
		temperatureGroup.GET("/hamster/:idHamster", r.GetByHamsterController.GetByHamster)
	}
}
