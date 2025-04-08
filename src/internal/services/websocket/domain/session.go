package domain

import (
	"log"
	"time"
	"github.com/gorilla/websocket"
)

type Session struct {
	conn         *websocket.Conn
	UserID       int32  
	Role         string 
	closeHandler func(userID int32)
}

func NewSession(conn *websocket.Conn, userID int32, role string) *Session {
	return &Session{
		conn:   conn,
		UserID: userID,
		Role:   role,
	}
}

func (s *Session) SetCloseHandler(handler func(userID int32)) {
	s.closeHandler = handler
}

func (s *Session) StartHandling() {
	defer func() {
		s.conn.Close()
		if s.closeHandler != nil {
			s.closeHandler(s.UserID)
		}
	}()

	for {

		if _, _, err := s.conn.ReadMessage(); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error: %v", err)
			}
			break
		}
		time.Sleep(17 * time.Millisecond)
	}
}

func (s *Session) SendUpdate(cageData interface{}) error {
	return s.conn.WriteJSON(cageData)
}