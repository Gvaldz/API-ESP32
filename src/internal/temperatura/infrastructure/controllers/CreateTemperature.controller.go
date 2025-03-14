package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"esp32/src/internal/temperatura/application"
	"esp32/src/internal/temperatura/domain"
)

type CreateTemperatureController struct {
	createTemperature *application.CreateTemperature
}

func NewCreateTemperatureController(createTemperature *application.CreateTemperature) *CreateTemperatureController {
	return &CreateTemperatureController{createTemperature: createTemperature}
}

func (h *CreateTemperatureController) Create(c *gin.Context) {
	var temperatureRequest domain.Temperature
	if err := c.ShouldBindJSON(&temperatureRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.createTemperature.Execute(temperatureRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Temperatura creada correctamente", "temperature": temperatureRequest})
}