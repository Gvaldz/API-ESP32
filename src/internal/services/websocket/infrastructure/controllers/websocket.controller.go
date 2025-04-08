package controllers

import (
	"esp32/src/core"
	"esp32/src/internal/services/websocket/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebSocketController struct {
	wsService    *application.WebSocketService
	tokenService core.JWTService
}

func NewWebSocketController(
	wsService *application.WebSocketService,
	tokenService core.JWTService,
) *WebSocketController {
	return &WebSocketController{
		wsService:    wsService,
		tokenService: tokenService,
	}
}

func (c *WebSocketController) HandleConnection(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token requerido"})
		return
	}

	userID, userType, err := c.tokenService.ValidateToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
		return
	}

	if err := c.wsService.HandleConnection(
		ctx.Writer,
		ctx.Request,
		userID,
		userType,
	); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "no se pudo establecer conexión WebSocket",
		})
	}
}
