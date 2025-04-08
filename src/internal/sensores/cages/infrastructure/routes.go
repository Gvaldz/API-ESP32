package infrastructure

import (
	"esp32/src/core"
	"esp32/src/internal/sensores/cages/infrastructure/controllers"
	tokenService "esp32/src/internal/services/auth/domain"
	"esp32/src/server/middleware"

	"github.com/gin-gonic/gin"
)

type CageRoutes struct {
	CreateCageController     *controllers.CreateCageController
	GetAllCagesController    *controllers.GetAllCagesController
	GetCageController        *controllers.GetCageByIDController
	GetCagesByUserController *controllers.GetCagesByUserController
	UpdateCageController     *controllers.UpdateCageController
	TokenService             tokenService.TokenService
	AuthRepo                 *core.AuthRepository
}

func NewCageRoutes(
	createCageController *controllers.CreateCageController,
	getAllCagesController *controllers.GetAllCagesController,
	getCageController *controllers.GetCageByIDController,
	getCagesByUserController *controllers.GetCagesByUserController,
	updateCageController *controllers.UpdateCageController,
	tokenService tokenService.TokenService,
	authRepo *core.AuthRepository,
) *CageRoutes {
	return &CageRoutes{
		CreateCageController:     createCageController,
		GetAllCagesController:    getAllCagesController,
		GetCageController:        getCageController,
		GetCagesByUserController: getCagesByUserController,
		UpdateCageController:     updateCageController,
		TokenService:             tokenService,
		AuthRepo:                 authRepo,
	}
}

func (r *CageRoutes) AttachRoutes(router *gin.Engine) {
	userAuth := middleware.AuthMiddleware(r.TokenService, r.AuthRepo, "usuario")
	adminAuth := middleware.AuthMiddleware(r.TokenService, r.AuthRepo, "administrador")

	userGroup := router.Group("/cages")
	userGroup.Use(userAuth)
	{
		userGroup.GET("/:id", r.GetCageController.GetCageByID)
		userGroup.GET("/user/:id", r.GetCagesByUserController.GetByUser)
		userGroup.PUT("/:id", r.UpdateCageController.UpdateUser)
	}

	adminGroup := router.Group("/admin/cages")
	adminGroup.Use(adminAuth)
	{
		adminGroup.POST("", r.CreateCageController.Create)
		adminGroup.GET("", r.GetAllCagesController.GetAllCages)
		adminGroup.PUT("/:id", r.UpdateCageController.UpdateUser)
	}
}
