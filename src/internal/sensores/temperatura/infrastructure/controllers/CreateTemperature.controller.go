package controllers

import (
	core "esp32/src/core"
	cages "esp32/src/internal/sensores/cages/domain"
	"esp32/src/internal/sensores/temperatura/application"
	"esp32/src/internal/sensores/temperatura/domain"
	websocket "esp32/src/internal/services/websocket/application"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateTemperatureController struct {
	createTemperature *application.CreateTemperature
	wsService         *websocket.WebSocketService
	cageRepo          cages.CageRepository
	userRepo          *core.UserRepository
}

func NewCreateTemperatureController(
	createTemperature *application.CreateTemperature,
	wsService *websocket.WebSocketService,
	cageRepo cages.CageRepository,
	userRepo *core.UserRepository,
) *CreateTemperatureController {
	return &CreateTemperatureController{
		createTemperature: createTemperature,
		wsService:         wsService,
		cageRepo:          cageRepo,
		userRepo:          userRepo,
	}
}

func (h *CreateTemperatureController) Create(c *gin.Context) {
	var temperatureRequest domain.Temperature
	if err := c.ShouldBindJSON(&temperatureRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Creando temperatura desde HTTP: %+v\n", temperatureRequest)

	err := h.createTemperature.Execute(temperatureRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Temperatura creada correctamente", "temperature": temperatureRequest})
}

func (h *CreateTemperatureController) ProcessTemperature(temperature domain.Temperature) error {
	log.Printf("[DEBUG] Iniciando procesamiento de temperatura: %+v", temperature)

	if err := h.createTemperature.Execute(temperature); err != nil {
		log.Printf("[ERROR] Fallo al guardar temperatura: %v", err)
		return err
	}
	log.Printf("[DEBUG] Temperatura guardada en BD: %+v", temperature)

	cage, err := h.cageRepo.GetCageByID(temperature.IDHamster)
	if err != nil {
		log.Printf("[ERROR] No se pudo obtener jaula %s: %v", temperature.IDHamster, err)
		return err
	}
	log.Printf("[DEBUG] Jaula obtenida: %+v", cage)

	wsData := gin.H{
		"event":     "new_temperature",
		"data":      temperature,
		"cage_id":   temperature.IDHamster,
		"timestamp": time.Now().Unix(),
	}

	if err := h.wsService.NotifyUser(cage.Idusuario, wsData); err != nil {
		log.Printf("[WARN] Error notificando usuario %d via WebSocket: %v", cage.Idusuario, err)
	} else {
		log.Printf("[DEBUG] Notificación WebSocket enviada al usuario %d", cage.Idusuario)
	}

	log.Printf("[INFO] Procesamiento completado para temperatura: %+v", temperature)
	return nil
}
