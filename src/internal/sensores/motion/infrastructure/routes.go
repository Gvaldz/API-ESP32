package infrastructure

import (
	"esp32/src/internal/sensores/motion/infrastructure/controllers"

	"github.com/gin-gonic/gin"
)

type MotionRoutes struct {
	CreateMotionController *controllers.CreateMotionController
	GetByHamsterController *controllers.GetByHamsterController
}

func NewMotionRoutes(
	createMotionCreateMotionController *controllers.CreateMotionController,
	getByHamsterController *controllers.GetByHamsterController,
) *MotionRoutes {
	return &MotionRoutes{
		CreateMotionController: createMotionCreateMotionController,
		GetByHamsterController: getByHamsterController,
	}
}

func (r *MotionRoutes) AttachRoutes(router *gin.Engine) {
	motionsGroup := router.Group("/motions")
	{
		motionsGroup.POST("", r.CreateMotionController.Create)
		motionsGroup.GET("/hamster/:idHamster", r.GetByHamsterController.GetByHamster)
	}
}
