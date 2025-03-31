package application

import (
	"esp32/src/internal/websocket/domain"
	"fmt"
)

// WebSocketNotifier gestiona las notificaciones WebSocket
type WebSocketNotifier struct {
	wsService domain.WebSocketService
}

// NewWebSocketNotifier crea un nuevo servicio de notificación
func NewWebSocketNotifier(wsService domain.WebSocketService) *WebSocketNotifier {
	return &WebSocketNotifier{wsService: wsService}
}

// NotifyAll envía una alerta a todos los clientes WebSocket conectados
func (n *WebSocketNotifier) NotifyAll(sensor string, value float64, timestamp string) {
	message := fmt.Sprintf("⚠️ Advertencia: %s fuera de rango: %.2f", sensor, value)

	wsMessage := domain.WebSocketMessage{
		Sensor:    sensor,
		Message:   message,
		Value:     value,
		Timestamp: timestamp,
	}

	n.wsService.BroadcastMessage(wsMessage)
}
