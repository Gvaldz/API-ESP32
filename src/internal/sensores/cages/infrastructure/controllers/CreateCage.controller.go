package controllers

import (
	"esp32/src/internal/sensores/cages/application"
	"esp32/src/internal/sensores/cages/domain"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateCageController struct {
	createCage *application.CreateCage
}

func NewCreateCageController(createCage *application.CreateCage) *CreateCageController {
	return &CreateCageController{
		createCage: createCage,
	}
}

func (h *CreateCageController) Create(c *gin.Context) {
	var cageRequest domain.Cage
	if err := c.ShouldBindJSON(&cageRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.createCage.Execute(cageRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "jaula creada correctamente", "cage": cageRequest})
}

func (h *CreateCageController) ProcessCage(cage domain.Cage) error {
	fmt.Printf("Procesando creación de jaula: %+v\n", cage)
	return h.createCage.Execute(cage)
}
