package controllers

import (
	"context"
	"esp32/src/core"
	cages "esp32/src/internal/sensores/cages/domain"
	"esp32/src/internal/sensores/humidity/application"
	"esp32/src/internal/sensores/humidity/domain"
	fcm "esp32/src/internal/services/fcm"
	websocket "esp32/src/internal/services/websocket/application"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateHumidityController struct {
	createHumidity *application.CreateHumidity
	wsService      *websocket.WebSocketService
	cageRepo       cages.CageRepository
	userRepo       *core.UserRepository
	fcmSender      *fcm.FCMSender
}

func NewCreateHumidityController(
	createHumidity *application.CreateHumidity,
	wsService *websocket.WebSocketService,
	cageRepo cages.CageRepository,
	userRepo *core.UserRepository,
	fcmSender *fcm.FCMSender,
) *CreateHumidityController {
	return &CreateHumidityController{
		createHumidity: createHumidity,
		wsService:      wsService,
		cageRepo:       cageRepo,
		userRepo:       userRepo,
		fcmSender:      fcmSender,
	}
}

func (h *CreateHumidityController) Create(c *gin.Context) {
	var humidityRequest domain.Humidity
	if err := c.ShouldBindJSON(&humidityRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Creando humedad: %+v\n", humidityRequest)

	err := h.createHumidity.Execute(humidityRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Humedad creada correctamente", "humidity": humidityRequest})
}

func (h *CreateHumidityController) ProcessHumidity(humidity domain.Humidity) error {
	log.Printf("[DEBUG] Iniciando procesamiento de humedad: %+v", humidity)
	
	if err := h.createHumidity.Execute(humidity); err != nil {
		log.Printf("[ERROR] Fallo al guardar humedad: %v", err)
		return err
	}
	log.Printf("[DEBUG] Humedad guardada en BD: %+v", humidity)

	cage, err := h.cageRepo.GetCageByID(humidity.IDHamster)
	if err != nil {
		log.Printf("[ERROR] No se pudo obtener jaula %s: %v", humidity.IDHamster, err)
		return err
	}
	log.Printf("[DEBUG] Jaula obtenida: %+v", cage)

	wsData := gin.H{
		"event":     "new_humidity",
		"data":      humidity,
		"cage_id":   humidity.IDHamster,
		"timestamp": time.Now().Unix(),
	}

	if err := h.wsService.NotifyUser(cage.Idusuario, wsData); err != nil {
		log.Printf("[WARN] Error notificando usuario %d via WebSocket: %v", cage.Idusuario, err)
	} else {
		log.Printf("[DEBUG] Notificación WebSocket enviada al usuario %d", cage.Idusuario)
	}

	user, err := h.userRepo.GetUserByID(cage.Idusuario)
	if err != nil {
		log.Printf("[ERROR] No se pudo obtener usuario %d: %v", cage.Idusuario, err)
		return fmt.Errorf("error obteniendo usuario: %v", err)
	}

	if user.FCMToken != "" {
		payload := fcm.NotificationPayload{
			Title: "Nueva humedad registrada",
			Body:  fmt.Sprintf("Jaula %s: %.2f%%", humidity.IDHamster, humidity.Humedad),
			Data: map[string]string{
				"cage_id":  fmt.Sprintf("%d", humidity.IDHamster),
				"humidity": fmt.Sprintf("%.2f", humidity.Humedad),
				"timestamp": time.Now().Format(time.RFC3339),
			},
		}

		if err := h.fcmSender.SendNotification(context.Background(), user.FCMToken, payload); err != nil {
			log.Printf("[ERROR] Fallo al enviar notificación FCM: %v", err)
		} else {
			log.Printf("[DEBUG] Notificación FCM enviada a token: %s", user.FCMToken)
		}
	} else {
		log.Printf("[DEBUG] Usuario %d no tiene FCMToken registrado", user.IdUsuario)
	}

	log.Printf("[INFO] Procesamiento completado para humedad: %+v", humidity)
	return nil
}