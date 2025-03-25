package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"esp32/src/internal/motion/application"
	"esp32/src/internal/motion/domain"
)

type CreateMotionController struct {
	createMotion *application.CreateMotion
}

func NewCreateMotionController(createMotion *application.CreateMotion) *CreateMotionController {
	return &CreateMotionController{createMotion: createMotion}
}

func (h *CreateMotionController) Create(c *gin.Context) {
	var motionRequest domain.Motion
	if err := c.ShouldBindJSON(&motionRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Creando movimiento desde HTTP: %+v\n", motionRequest)

	err := h.createMotion.Execute(motionRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Movimiento creado correctamente", "movimiento": motionRequest})
}

func (h *CreateMotionController) ProcessMotion(motion domain.Motion) error {
	fmt.Printf("Procesando movivimiento desde AMQP: %+v\n", motion)
	return h.createMotion.Execute(motion)
}
