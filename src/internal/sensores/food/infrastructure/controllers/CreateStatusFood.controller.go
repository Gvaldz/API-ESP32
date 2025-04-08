package controllers

import (
	"context"
	"esp32/src/core"
	cages "esp32/src/internal/sensores/cages/domain"
	"esp32/src/internal/sensores/food/application"
	"esp32/src/internal/sensores/food/domain"
	fcm "esp32/src/internal/services/fcm"
	websocket "esp32/src/internal/services/websocket/application"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateStatusFoodController struct {
	createStatusFood *application.CreateStatusFood
	wsService        *websocket.WebSocketService
	cageRepo         cages.CageRepository
	userRepo         *core.UserRepository
	fcmSender        *fcm.FCMSender
}

func NewCreateStatusFoodController(
	createStatusFood *application.CreateStatusFood,
	wsService *websocket.WebSocketService,
	cageRepo cages.CageRepository,
	userRepo *core.UserRepository,
	fcmSender *fcm.FCMSender,
) *CreateStatusFoodController {
	return &CreateStatusFoodController{
		createStatusFood: createStatusFood,
		wsService:        wsService,
		cageRepo:         cageRepo,
		userRepo:         userRepo,
		fcmSender:        fcmSender,
	}
}

func (h *CreateStatusFoodController) Create(c *gin.Context) {
	var foodRequest domain.Food
	if err := c.ShouldBindJSON(&foodRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Creando estatus de alimento: %+v\n", foodRequest)

	err := h.createStatusFood.Execute(foodRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Estatus de alimento creado correctamente", "food": foodRequest})
}

func (h *CreateStatusFoodController) ProcessFood(food domain.Food) error {
    log.Printf("[DEBUG] Iniciando procesamiento de alimento: %+v", food)
    
    if err := h.createStatusFood.Execute(food); err != nil {
        log.Printf("[ERROR] Fallo al guardar estatus de alimento: %v", err)
        return err
    }
    
    cage, err := h.cageRepo.GetCageByID(food.IDHamster)
    if err != nil {
        log.Printf("[ERROR] No se pudo obtener jaula %s: %v", food.IDHamster, err)
        return err
    }

    wsData := gin.H{
        "cage_id": food.IDHamster,
        "data": gin.H{
            "idalimento":    food.IDalimento,
            "idhamster":     food.IDHamster,
            "alimento":      food.Alimento,
            "porcentaje":    food.Porcentaje,
            "hora_registro": food.HoraRegistro,
        },
        "event":     "new_food",
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
        status := "sin alimento"
        if food.Alimento == 1 {
            status = fmt.Sprintf("con alimento (%.1f%%)", food.Porcentaje)
        }
        
        payload := fcm.NotificationPayload{
            Title: "Estado del alimento",
            Body:  fmt.Sprintf("Jaula %s: %s", food.IDHamster, status),
            Data: map[string]string{
                "cage_id":      food.IDHamster,
                "alimento":     fmt.Sprintf("%d", food.Alimento),
                "porcentaje":   fmt.Sprintf("%.1f", food.Porcentaje),
                "timestamp":    time.Now().Format(time.RFC3339),
                "event":        "new_food",
            },
        }

        if err := h.fcmSender.SendNotification(context.Background(), user.FCMToken, payload); err != nil {
            log.Printf("[ERROR] Fallo al enviar notificación FCM: %v", err)
        }
    }

    log.Printf("[INFO] Notificación de alimento enviada: %+v", wsData)
    return nil
}