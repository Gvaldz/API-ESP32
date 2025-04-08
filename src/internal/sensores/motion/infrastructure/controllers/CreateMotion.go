package controllers

import (
	"context"
	"esp32/src/core"
	cages "esp32/src/internal/sensores/cages/domain"
	"esp32/src/internal/sensores/motion/application"
	"esp32/src/internal/sensores/motion/domain"
	fcm "esp32/src/internal/services/fcm"
	websocket "esp32/src/internal/services/websocket/application"
	"fmt"
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
	fcmSender    *fcm.FCMSender
}

func NewCreateMotionController(
	createMotion *application.CreateMotion,
	wsService *websocket.WebSocketService,
	cageRepo cages.CageRepository,
	userRepo *core.UserRepository,
	fcmSender *fcm.FCMSender,
) *CreateMotionController {
	return &CreateMotionController{
		createMotion: createMotion,
		wsService:    wsService,
		cageRepo:     cageRepo,
		userRepo:     userRepo,
		fcmSender:    fcmSender,
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
            "movimiento":   motion.Movimiento,
            "hora_registro": motion.HoraRegistro,
        },
        "event":     "new_motion",
        "timestamp": time.Now().Unix(),
    }

    if err := h.wsService.NotifyUser(cage.Idusuario, wsData); err != nil {
        log.Printf("[WARN] Error notificando usuario %d via WebSocket: %v", cage.Idusuario, err)
    }

    user, err := h.userRepo.GetUserByID(cage.Idusuario)
    if err != nil {
        log.Printf("[ERROR] No se pudo obtener usuario %d: %v", cage.Idusuario, err)
        return fmt.Errorf("error obteniendo usuario: %v", err)
    }

    if user.FCMToken != "" {
        status := "sin movimiento"
        if motion.Movimiento {
            status = "movimiento detectado"
        }
        
        payload := fcm.NotificationPayload{
            Title: "Detección de movimiento",
            Body:  fmt.Sprintf("Jaula %s: %s", motion.IDHamster, status),
            Data: map[string]string{
                "cage_id":    motion.IDHamster,
                "movimiento": fmt.Sprintf("%t", motion.Movimiento),
                "timestamp":  time.Now().Format(time.RFC3339),
                "event":      "new_motion",
            },
        }

        if err := h.fcmSender.SendNotification(context.Background(), user.FCMToken, payload); err != nil {
            log.Printf("[ERROR] Fallo al enviar notificación FCM: %v", err)
        }
    }

    log.Printf("[INFO] Notificación de movimiento enviada: %+v", wsData)
    return nil
}