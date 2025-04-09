package controllers

import (
	"esp32/src/core"
	cages "esp32/src/internal/sensores/cages/domain"
	"esp32/src/internal/sensores/motion/application"
	"esp32/src/internal/sensores/motion/domain"
	websocket "esp32/src/internal/services/websocket/application"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateMotionController struct {
	createMotion *application.CreateMotion
	wsService    *websocket.WebSocketService
	cageRepo     cages.CageRepository
	userRepo     *core.UserRepository
}

func NewCreateMotionController(
	createMotion *application.CreateMotion,
	wsService *websocket.WebSocketService,
	cageRepo cages.CageRepository,
	userRepo *core.UserRepository,
) *CreateMotionController {
	return &CreateMotionController{
		createMotion: createMotion,
		wsService:    wsService,
		cageRepo:     cageRepo,
		userRepo:     userRepo,
	}
}

func (h *CreateMotionController) Create(c *gin.Context) {
	var motionRequest domain.Motion
	if err := c.ShouldBindJSON(&motionRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Creando movimiento: %+v\n", motionRequest)

	err := h.createMotion.Execute(motionRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Movimiento creado correctamente", "motion": motionRequest})
}

func (h *CreateMotionController) ProcessMotion(motion domain.Motion) error {
	log.Printf("[DEBUG] Iniciando procesamiento de movimiento: %+v", motion)

	if err := h.createMotion.Execute(motion); err != nil {
		log.Printf("[ERROR] Fallo al guardar movimiento: %v", err)
		return err
	}

	cage, err := h.cageRepo.GetCageByID(motion.IDHamster)
	if err != nil {
		log.Printf("[ERROR] No se pudo obtener jaula %s: %v", motion.IDHamster, err)
		return err
	}

	wsData := gin.H{
		"cage_id": motion.IDHamster,
		"data": gin.H{
			"idhamster":     motion.IDHamster,
			"movimiento":    motion.Movimiento,
			"hora_registro": motion.HoraRegistro,
		},
		"event":     "new_motion",
		"timestamp": time.Now().Unix(),
	}

	if err := h.wsService.NotifyUser(cage.Idusuario, wsData); err != nil {
		log.Printf("[WARN] Error notificando usuario %d via WebSocket: %v", cage.Idusuario, err)
	}

	log.Printf("[INFO] Notificación de movimiento enviada: %+v", wsData)
	return nil
}
