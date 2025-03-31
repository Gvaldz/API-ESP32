package controllers

import (
	"esp32/src/internal/cages/application"
	"esp32/src/internal/cages/domain"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateCageController struct {
	updateCageController *application.UpdateCage
}

func NewUpdateCageController(updateCageController *application.UpdateCage) *UpdateCageController {
	return &UpdateCageController{
		updateCageController: updateCageController,
	}
}

func (c *UpdateCageController) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	
	var user domain.Cage
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	var idInt int32
	if _, err := fmt.Sscanf(id, "%d", &idInt); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	
	if err := c.updateCageController.Execute(idInt, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "jaula actualizada correctamente"})
}