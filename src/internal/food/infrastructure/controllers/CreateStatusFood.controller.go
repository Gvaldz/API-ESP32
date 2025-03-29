package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"esp32/src/internal/food/application"
	"esp32/src/internal/food/domain"
)

type CreateStatusFoodController struct {
	createStatusFood *application.CreateStatusFood
}

func NewCreateStatusFoodController(createStatusFood *application.CreateStatusFood) *CreateStatusFoodController {
	return &CreateStatusFoodController{createStatusFood: createStatusFood}
}

func (h *CreateStatusFoodController) Create(c *gin.Context) {
	var foodRequest domain.Food
	if err := c.ShouldBindJSON(&foodRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Creando estatus de alimento desde HTTP: %+v\n", foodRequest)

	err := h.createStatusFood.Execute(foodRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "estatus de alimento creada correctamente", "food": foodRequest})
}

func (h *CreateStatusFoodController) ProcessFood(food domain.Food) error {
	fmt.Printf("Procesando estatus de alimento desde AMQP: %+v\n", food)
	return h.createStatusFood.Execute(food)
}
