package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"esp32/src/internal/motion/application"
	"esp32/src/internal/motion/domain"
	cages 		"esp32/src/internal/cages/domain"
	websocket   "esp32/src/internal/websocket/application"

)

type CreateMotionController struct {
	createMotion *application.CreateMotion
    wsService    *websocket.WebSocketService
    cageRepo     cages.CageRepository
}

func NewCreateMotionController(
	createMotion *application.CreateMotion,
	wsService	 *websocket.WebSocketService,
    cageRepo 	 cages.CageRepository,
	) *CreateMotionController {
	return &CreateMotionController{        
		createMotion: createMotion,
        wsService:   wsService,
        cageRepo:    cageRepo,}
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
    fmt.Printf("Procesando movimiento desde AMQP: %+v\n", motion)
    
    if err := h.createMotion.Execute(motion); err != nil {
        return err
    }
    
    cage, err := h.cageRepo.GetCageByID(motion.IDHamster)
    if err != nil {
        return err
    }
    
    if err := h.wsService.NotifyUser(cage.Idusuario, gin.H{
        "event": "new_motion",
        "data":  motion,
        "cage_id": motion.IDHamster,
        "timestamp": time.Now().Unix(),
    }); err != nil {
        log.Printf("Error notificando al usuario %d: %v", cage.Idusuario, err)
    }
    
    return nil
}