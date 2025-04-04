package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	users "esp32/src/internal/users/infrastructure"
	fcm			 "esp32/src/internal/fcm"
	cages "esp32/src/internal/cages/domain"
	"esp32/src/internal/temperatura/application"
	"esp32/src/internal/temperatura/domain"
	websocket "esp32/src/internal/websocket/application"

	"github.com/gin-gonic/gin"
)
type CreateTemperatureController struct {
	createTemperature 	*application.CreateTemperature
    wsService    		*websocket.WebSocketService
    cageRepo    		cages.CageRepository
	userRepo    		users.UsersRepo
	fcmSender     		*fcm.FCMSender
}

func NewCreateTemperatureController(
	createTemperature *application.CreateTemperature,
	wsService	 *websocket.WebSocketService,
    cageRepo 	 cages.CageRepository,
	userRepo	 users.UsersRepo,
	fcmSender	 *fcm.FCMSender,
	) *CreateTemperatureController {
	return &CreateTemperatureController{        
		createTemperature: createTemperature,
        wsService:   wsService,
        cageRepo:    cageRepo,
		fcmSender:        fcmSender, 
}
}

func (h *CreateTemperatureController) Create(c *gin.Context) {
	var temperatureRequest domain.Temperature
	if err := c.ShouldBindJSON(&temperatureRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Creando temperatura desde HTTP: %+v\n", temperatureRequest)

	err := h.createTemperature.Execute(temperatureRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Temperatura creada correctamente", "temperature": temperatureRequest})
}

func (h *CreateTemperatureController) ProcessTemperature(temperature domain.Temperature) error {
    fmt.Printf("Procesando temperatura desde AMQP: %+v\n", temperature)
    
    if err := h.createTemperature.Execute(temperature); err != nil {
        return err
    }
    
    cage, err := h.cageRepo.GetCageByID(temperature.IDHamster)
    if err != nil {
        return err
    }
    
    if err := h.wsService.NotifyUser(cage.Idusuario, gin.H{
        "event": "new_temperature",
        "data":  temperature,
        "cage_id": temperature.IDHamster,
        "timestamp": time.Now().Unix(),
    }); err != nil {
        log.Printf("Error notificando al usuario %d via WebSocket: %v", cage.Idusuario, err)
    }
    
    /*user, err := h.userRepo.GetUserByID(cage.Idusuario)
    if err != nil {
        return fmt.Errorf("error obteniendo usuario: %v", err)
    }
    
    if user.FCMToken != "" {
        payload := fcm.NotificationPayload{
            Title: "Nueva temperatura registrada",
            Body:  fmt.Sprintf("Jaula %d: %.2f°C", temperature.IDHamster, temperature.Temperatura),
            Data: map[string]string{
                "cage_id":    (string(temperature.IDHamster)),
                "temperature": fmt.Sprintf("%.2f", temperature.Temperatura),
                "timestamp":  time.Now().Format(time.RFC3339),
            },
        }
        
        if err := h.fcmSender.SendNotification(context.Background(), user.FCMToken, payload); err != nil {
            log.Printf("Error enviando notificación FCM: %v", err)
        }
    }
    */
    return nil 
}