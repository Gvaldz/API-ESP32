package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"esp32/src/internal/temperatura/application"
	"esp32/src/internal/temperatura/domain"
)

type CreateTemperatureController struct {
	createTemperature *application.CreateTemperature
}

// Constructor del controlador
func NewCreateTemperatureController(createTemperature *application.CreateTemperature) *CreateTemperatureController {
	return &CreateTemperatureController{createTemperature: createTemperature}
}

// Método HTTP para crear temperatura
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

// Método que se usará en el Consumer AMQP
func (h *CreateTemperatureController) ProcessTemperature(temperature domain.Temperature) error {
	fmt.Printf("Procesando temperatura desde AMQP: %+v\n", temperature)
	return h.createTemperature.Execute(temperature)
}
