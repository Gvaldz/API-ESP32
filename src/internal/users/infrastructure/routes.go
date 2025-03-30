package infrastructure

import (
	"esp32/src/internal/users/infrastructure/controllers"
	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	CreateUserController       *controllers.CreateUserController
	GetUserController          *controllers.GetByUserIDController
	UpdateUserController       *controllers.UpdateUserController
	UpdatePasswordController   *controllers.UpdatePasswordController
}

func NewUserRoutes(
	createUserController 	   *controllers.CreateUserController,
	getUserController 		   *controllers.GetByUserIDController,
	updateUserController 	   *controllers.UpdateUserController,
	updatePasswordController   *controllers.UpdatePasswordController,
) *UserRoutes {
	return &UserRoutes{
		CreateUserController:     createUserController,
		GetUserController:        getUserController,
		UpdateUserController:     updateUserController,
		UpdatePasswordController: updatePasswordController,
	}
}

func (r *UserRoutes) AttachRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("", r.CreateUserController.CreateUser)
		userGroup.GET("/:id", r.GetUserController.GetByUserID)
		userGroup.PUT("/:id", r.UpdateUserController.UpdateUser)
		userGroup.PUT("/:id/password", r.UpdatePasswordController.UpdatePassword)
	}
}