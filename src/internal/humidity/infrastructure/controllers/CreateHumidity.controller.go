package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	cages "esp32/src/internal/cages/domain"
	"esp32/src/internal/humidity/application"
	"esp32/src/internal/humidity/domain"
	websocket "esp32/src/internal/websocket/application"

	"github.com/gin-gonic/gin"
)

type CreateHumidityController struct {
	createHumidity *application.CreateHumidity
	wsService    *websocket.WebSocketService
    cageRepo     cages.CageRepository
}

func NewCreateHumidityController(
	createHumidity *application.CreateHumidity,
	wsService	 *websocket.WebSocketService,
    cageRepo 	 cages.CageRepository,
	) *CreateHumidityController {
	return &CreateHumidityController{        
		createHumidity: createHumidity,
        wsService:   wsService,
        cageRepo:    cageRepo,}
}

func (h *CreateHumidityController) Create(c *gin.Context) {
	var humidityRequest domain.Humidity
	if err := c.ShouldBindJSON(&humidityRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Creando humedad: %+v\n", humidityRequest)

	err := h.createHumidity.Execute(humidityRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "humedad creada correctamente", "humidity": humidityRequest})
}

func (h *CreateHumidityController) ProcessHumidity(humidity domain.Humidity) error {
    fmt.Printf("Procesando temperatura desde AMQP: %+v\n", humidity)
    
    if err := h.createHumidity.Execute(humidity); err != nil {
        return err
    }
    
    cage, err := h.cageRepo.GetCageByID(humidity.IDHamster)
    if err != nil {
        return err
    }
    
    if err := h.wsService.NotifyUser(cage.Idusuario, gin.H{
        "event": "new_humidity",
        "data":  humidity,
        "cage_id": humidity.IDHamster,
        "timestamp": time.Now().Unix(),
    }); err != nil {
        log.Printf("Error notificando al usuario %d: %v", cage.Idusuario, err)
    }
    
    return nil
}