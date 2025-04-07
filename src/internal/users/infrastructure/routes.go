package infrastructure

import (
    "esp32/src/internal/users/infrastructure/controllers"
    "github.com/gin-gonic/gin"
)

type UserRoutes struct {
    CreateUserController       *controllers.CreateUserController
    GetAllUsersController      *controllers.GetAllUsersController
    GetUserController          *controllers.GetByUserIDController
    UpdateUserController       *controllers.UpdateUserController
    UpdatePasswordController   *controllers.UpdatePasswordController
    DeleteUserController       *controllers.DeleteUserController
    FCMController              *controllers.FCMController
    AuthMiddleware             gin.HandlerFunc 
}

func NewUserRoutes(
    createUserController      *controllers.CreateUserController,
    getAllUsersController     *controllers.GetAllUsersController,
    getUserController         *controllers.GetByUserIDController,
    updateUserController      *controllers.UpdateUserController,
    updatePasswordController *controllers.UpdatePasswordController,
    deleteUserController      *controllers.DeleteUserController,
    fcmController            *controllers.FCMController,
    authMiddleware           gin.HandlerFunc, 
) *UserRoutes {
    return &UserRoutes{
        CreateUserController:     createUserController,
        GetAllUsersController:    getAllUsersController,
        GetUserController:        getUserController,
        UpdateUserController:    updateUserController,
        UpdatePasswordController: updatePasswordController,
        DeleteUserController:    deleteUserController,
        FCMController:           fcmController,
        AuthMiddleware:          authMiddleware, 
    }
}

func (r *UserRoutes) AttachRoutes(router *gin.Engine) {
    userGroup := router.Group("/users")
    {
        userGroup.POST("", r.CreateUserController.CreateUser)
        userGroup.GET("", r.GetAllUsersController.GetAll)
        userGroup.GET("/:id", r.GetUserController.GetByUserID)
        userGroup.PUT("/:id", r.UpdateUserController.UpdateUser)
        userGroup.PUT("/password/:id", r.UpdatePasswordController.UpdatePassword)
        userGroup.DELETE("/:id", r.DeleteUserController.Delete)
        userGroup.POST("/fcm-token", r.AuthMiddleware, r.FCMController.RegisterToken)
    }
}