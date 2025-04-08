package controllers

import (
	"esp32/src/internal/sensores/cages/application"
	"esp32/src/internal/sensores/cages/domain"
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

	if err := c.updateCageController.Execute(id, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "jaula actualizada correctamente"})
}
