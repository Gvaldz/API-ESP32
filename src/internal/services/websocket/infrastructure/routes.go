package infrastructure

import (
	"esp32/src/internal/services/websocket/infrastructure/controllers"

	"github.com/gin-gonic/gin"
)

type WebSocketRoutes struct {
	wsHandler *controllers.WebSocketController
}

func NewWebSocketRoutes(wsHandler *controllers.WebSocketController) *WebSocketRoutes {
	return &WebSocketRoutes{
		wsHandler: wsHandler,
	}
}

func (r *WebSocketRoutes) AttachRoutes(router *gin.Engine) {
	wsGroup := router.Group("/ws")
	{
		wsGroup.GET("/connect", r.wsHandler.HandleConnection)
	}
}
