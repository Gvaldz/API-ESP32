package controllers

import (
	"esp32/src/internal/users/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FCMController struct {
	userRepo domain.UserRepository
}

func NewFCMController(userRepo domain.UserRepository) *FCMController {
	return &FCMController{userRepo: userRepo}
}
func (c *FCMController) RegisterToken(ctx *gin.Context) {
    // Cambiar "user_id" por "userID" para que coincida con el middleware
    userIDInterface, exists := ctx.Get("userID")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
        return
    }

    userID, ok := userIDInterface.(int32) // El tipo correcto que usa el middleware
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "formato de ID inválido"})
        return
    }

    var request struct {
        Token string `json:"token"`
    }

    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "petición inválida"})
        return
    }

    // Convertir correctamente el int32 a string
    if err := c.userRepo.UpdateFCMToken(strconv.Itoa(int(userID)), request.Token); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar token"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "token actualizado"})
}