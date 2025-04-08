package domain

// WebSocketMessage representa la estructura de los mensajes que se enviarán a través de WebSocket
type WebSocketMessage struct {
	Sensor   string  `json:"sensor"`   // El nombre del sensor (ej. "temperatura")
	Message  string  `json:"message"`  // El mensaje a enviar (ej. "Temperatura fuera de rango")
	Value    float64 `json:"value"`    // El valor del sensor que está generando la alerta
	Timestamp string `json:"timestamp"` // La hora en que se generó la alerta
}
