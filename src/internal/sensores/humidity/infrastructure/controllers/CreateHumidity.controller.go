package controllers

import (
	"esp32/src/core"
	cages "esp32/src/internal/sensores/cages/domain"
	"esp32/src/internal/sensores/humidity/application"
	"esp32/src/internal/sensores/humidity/domain"
	websocket "esp32/src/internal/services/websocket/application"
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
}

func NewCreateHumidityController(
	createHumidity *application.CreateHumidity,
	wsService *websocket.WebSocketService,
	cageRepo cages.CageRepository,
	userRepo *core.UserRepository,
) *CreateHumidityController {
	return &CreateHumidityController{
		createHumidity: createHumidity,
		wsService:      wsService,
		cageRepo:       cageRepo,
		userRepo:       userRepo,
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

	log.Printf("[INFO] Procesamiento completado para humedad: %+v", humidity)
	return nil
}
