package controllers

import (
	"esp32/src/internal/users/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FCMController struct {
	userRepo domain.UserRepository
}

func NewFCMController(userRepo domain.UserRepository) *FCMController {
	return &FCMController{userRepo: userRepo}
}

func (c *FCMController) RegisterToken(ctx *gin.Context) {
    userIDInterface, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
        return
    }

    userID, ok := userIDInterface.(int32)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "formato de ID de usuario inválido"})
        return
    }

    var request struct {
        Token string `json:"token"`
    }

    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    if err := c.userRepo.UpdateFCMToken((string(userID)), request.Token); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update token"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "token updated"})
}