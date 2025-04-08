package infrastructure

import (
	"esp32/src/core"
	"esp32/src/internal/services/websocket/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebSocketHandler struct {
	wsService    *application.WebSocketService
	tokenService core.JWTService
}

func NewWebSocketHandler(wsService *application.WebSocketService, tokenService core.JWTService) *WebSocketHandler {
	return &WebSocketHandler{
		wsService:    wsService,
		tokenService: tokenService,
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token requerido"})
		return
	}

	userID, userType, err := h.tokenService.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
		return
	}

	if err := h.wsService.HandleConnection(c.Writer, c.Request, userID, userType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo establecer conexión WebSocket"})
	}
}
