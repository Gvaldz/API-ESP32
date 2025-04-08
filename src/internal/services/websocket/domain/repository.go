package domain

import "net/http"

// WebSocketService define la interfaz para enviar mensajes
type WebSocketService interface {
	BroadcastMessage(message WebSocketMessage)
	HandleConnection(w http.ResponseWriter, r *http.Request)
}