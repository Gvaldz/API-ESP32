package controllers

import (
	"esp32/src/internal/websocket/infrastructure"
	"github.com/gin-gonic/gin"
)

type WebSocketController struct {
	wsServer *infrastructure.WebSocketServer
}

func NewWebSocketController(wsServer *infrastructure.WebSocketServer) *WebSocketController {
	return &WebSocketController{wsServer: wsServer}
}

// Expone el WebSocket en la ruta `/ws`
func (w *WebSocketController) ConnectWebSocket(c *gin.Context) {
	w.wsServer.HandleConnection(c.Writer, c.Request)
}
