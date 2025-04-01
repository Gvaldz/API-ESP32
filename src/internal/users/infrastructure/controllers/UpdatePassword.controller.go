package controllers

import (
	"esp32/src/internal/users/application"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdatePasswordController struct {
	updatePasswordUC *application.UpdatePassword
}

func NewUpdatePasswordController(updatePasswordUC *application.UpdatePassword) *UpdatePasswordController {
	return &UpdatePasswordController{
		updatePasswordUC: updatePasswordUC,
	}
}

func (c *UpdatePasswordController) UpdatePassword(ctx *gin.Context) {
	id := ctx.Param("id")
	
	var request struct {
		NewPassword string `json:"newPassword" binding:"required"`
	}
	
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	var idInt int32
	if _, err := fmt.Sscanf(id, "%d", &idInt); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	
	if err := c.updatePasswordUC.Execute(idInt, request.NewPassword); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "Contraseña actualizada correctamente"})
}