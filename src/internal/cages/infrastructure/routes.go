package infrastructure

import (
    "esp32/src/server/middleware"
    "esp32/src/internal/cages/infrastructure/controllers"
    "github.com/gin-gonic/gin"
	authRepo     "esp32/src/internal/auth/infrastructure"
    tokenService "esp32/src/internal/auth/domain"
)

type CageRoutes struct {
    CreateCageController       *controllers.CreateCageController
    GetAllCagesController      *controllers.GetAllCagesController
    GetCageController          *controllers.GetCageByIDController
    GetCagesByUserController   *controllers.GetCagesByUserController
    UpdateCageController       *controllers.UpdateCageController
    TokenService                tokenService.TokenService 
    AuthRepo                   *authRepo.AuthRepositoryImpl
}

func NewCageRoutes(
    createCageController      *controllers.CreateCageController,
    getAllCagesController     *controllers.GetAllCagesController,
    getCageController         *controllers.GetCageByIDController,
    getCagesByUserController  *controllers.GetCagesByUserController,
    updateCageController      *controllers.UpdateCageController,
    tokenService             tokenService.TokenService,
    authRepo                 *authRepo.AuthRepositoryImpl,
) *CageRoutes {
    return &CageRoutes{
        CreateCageController:     createCageController,
        GetAllCagesController:    getAllCagesController,
        GetCageController:        getCageController,
        GetCagesByUserController: getCagesByUserController,
        UpdateCageController:     updateCageController,
        TokenService:            tokenService,
        AuthRepo:                authRepo,
    }
}

func (r *CageRoutes) AttachRoutes(router *gin.Engine) {
    userAuth := middleware.AuthMiddleware(r.TokenService, r.AuthRepo, "usuario")
    adminAuth := middleware.AuthMiddleware(r.TokenService, r.AuthRepo, "administrador")
    authGroup := router.Group("/cages")
    authGroup.Use(userAuth)
    
    {
        authGroup.GET("/my-cages", r.GetCagesByUserController.GetByUser) 
                adminGroup := authGroup.Group("")
        adminGroup.Use(adminAuth)
        {
            adminGroup.POST("", r.CreateCageController.Create)
            adminGroup.GET("", r.GetAllCagesController.GetAllCages)
            adminGroup.PUT("/:id", r.UpdateCageController.UpdateUser)
        }
        authGroup.GET("/:id", r.GetCageController.GetCageByID)
    }
}