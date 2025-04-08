package infrastructure

import (
	"esp32/src/internal/services/auth/infrastructure/controllers"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	LoginController *controllers.LoginController
}

func NewAuthRoutes(loginController *controllers.LoginController) *AuthRoutes {
	return &AuthRoutes{
		LoginController: loginController,
	}
}

func (r *AuthRoutes) AttachRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", r.LoginController.Login)
	}
}
