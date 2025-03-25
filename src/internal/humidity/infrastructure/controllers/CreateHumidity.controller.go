package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"esp32/src/internal/humidity/application"
	"esp32/src/internal/humidity/domain"
)

type CreateHumidityController struct {
	createHumidity *application.CreateHumidity
}

func NewCreateHumidityController(createHumidity *application.CreateHumidity) *CreateHumidityController {
	return &CreateHumidityController{createHumidity: createHumidity}
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
	fmt.Printf("Procesando humedad desde AMQP: %+v\n", humidity)
	return h.createHumidity.Execute(humidity)
}
