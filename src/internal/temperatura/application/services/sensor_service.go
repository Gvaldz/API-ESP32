package application

// import (
//     "esp32/src/core"
//     "esp32/src/internal/temperatura/domain"
//     "esp32/src/internal/temperatura/infrastructure/websocket"
// )

// type SensorService struct {
//     repo *core.MySQLRepository
//     ws   *websocket.WSServer
// }

// func NewSensorService(repo *core.MySQLRepository, ws *websocket.WSServer) *SensorService {
//     return &SensorService{repo: repo, ws: ws}
// }

// func (s *SensorService) ProcessSensorData(payload []byte) {
//     // Deserializar los datos y procesar la lógica
//     var data domain.SensorData
//     // Procesar datos del payload...

//     // Guardar datos en MySQL
//     err := s.repo.SaveData(data)
//     if err != nil {
//         // Manejar error
//         return
//     }

//     // Enviar los datos a los clientes WebSocket
//     s.ws.Broadcast(payload)
// }
