package infrastructure

import (
	"encoding/json"
	"esp32/src/internal/services/websocket/domain"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketServer gestiona conexiones WebSocket
type WebSocketServer struct {
	clients  map[*websocket.Conn]bool
	mutex    sync.Mutex
	upgrader websocket.Upgrader
}

// Constructor del servidor WebSocket
func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		clients: make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

// Envía un mensaje a todos los clientes WebSocket conectados
func (s *WebSocketServer) BroadcastMessage(message domain.WebSocketMessage) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	msgJSON, _ := json.Marshal(message)

	for client := range s.clients {
		err := client.WriteMessage(websocket.TextMessage, msgJSON)
		if err != nil {
			fmt.Println("Error al enviar mensaje:", err)
			client.Close()
			delete(s.clients, client)
		}
	}
}

// Manejo de nuevas conexiones WebSocket
func (s *WebSocketServer) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error al conectar WebSocket:", err)
		return
	}

	s.mutex.Lock()
	s.clients[conn] = true
	s.mutex.Unlock()

	fmt.Println("Nuevo cliente WebSocket conectado")

	// Leer mensajes y manejar desconexión
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			s.mutex.Lock()
			delete(s.clients, conn)
			s.mutex.Unlock()
			fmt.Println("Cliente WebSocket desconectado:", err)
			break
		}
	}
}
