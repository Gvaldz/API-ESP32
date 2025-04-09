package infrastructure

import (
	"esp32/src/internal/sensores/cages/infrastructure/controllers"
	"github.com/gin-gonic/gin"
)

type CageRoutes struct {
	CreateCageController     *controllers.CreateCageController
	GetAllCagesController    *controllers.GetAllCagesController
	GetCageController        *controllers.GetCageByIDController
	GetCagesByUserController *controllers.GetCagesByUserController
	UpdateCageController     *controllers.UpdateCageController
}

func NewCageRoutes(
	createCageController *controllers.CreateCageController,
	getAllCagesController *controllers.GetAllCagesController,
	getCageController *controllers.GetCageByIDController,
	getCagesByUserController *controllers.GetCagesByUserController,
	updateCageController *controllers.UpdateCageController,
) *CageRoutes {
	return &CageRoutes{
		CreateCageController:     createCageController,
		GetAllCagesController:    getAllCagesController,
		GetCageController:        getCageController,
		GetCagesByUserController: getCagesByUserController,
		UpdateCageController:     updateCageController,
	}
}

func (r *CageRoutes) AttachRoutes(router *gin.Engine) {
	cageGroup := router.Group("/cages")
	{
		cageGroup.GET("/:id", r.GetCageController.GetCageByID)
		cageGroup.GET("/user/:id", r.GetCagesByUserController.GetByUser)
		cageGroup.PUT("/:id", r.UpdateCageController.UpdateUser)
		cageGroup.POST("", r.CreateCageController.Create)
		cageGroup.GET("", r.GetAllCagesController.GetAllCages)
	}
}