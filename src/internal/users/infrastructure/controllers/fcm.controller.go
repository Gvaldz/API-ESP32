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
    userIDInterface, exists := ctx.Get("userID")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
        return
    }

    userID, ok := userIDInterface.(int32) 
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "formato de ID inválido"})
        return
    }

    var request struct {
        FCMToken string `json:"fcmToken" binding:"required"` 
    }

    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "petición inválida",
            "details": err.Error(), 
        })
        return
    }

    if request.FCMToken == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "fcmToken es requerido"})
        return
    }

    if err := c.userRepo.UpdateFCMToken(strconv.Itoa(int(userID)), request.FCMToken); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "error al actualizar token FCM",
            "details": err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "token FCM actualizado correctamente",
        "success": true,
    })
}