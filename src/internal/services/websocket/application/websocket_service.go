package application

import (
	"esp32/src/internal/services/websocket/domain"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketService struct {
	upgrader  websocket.Upgrader
	clients   map[int32]*domain.Session
	clientsMu sync.RWMutex
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clients: make(map[int32]*domain.Session),
	}
}

func (ws *WebSocketService) HandleConnection(w http.ResponseWriter, r *http.Request, userID int32, role string) error {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	session := domain.NewSession(conn, userID, role)

	ws.clientsMu.Lock()
	ws.clients[userID] = session
	ws.clientsMu.Unlock()

	session.SetCloseHandler(ws.removeClient)
	go session.StartHandling()

	return nil
}

func (ws *WebSocketService) removeClient(userID int32) {
	ws.clientsMu.Lock()
	defer ws.clientsMu.Unlock()

	delete(ws.clients, userID)
	log.Printf("Client disconnected: UserID %d", userID)
}

func (ws *WebSocketService) NotifyUser(userID int32, data interface{}) error {
	ws.clientsMu.RLock()
	defer ws.clientsMu.RUnlock()

	if client, ok := ws.clients[userID]; ok {
		return client.SendUpdate(data)
	}
	return nil
}
