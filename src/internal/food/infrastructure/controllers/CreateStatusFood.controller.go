package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	cages "esp32/src/internal/cages/domain"
	"esp32/src/internal/food/application"
	"esp32/src/internal/food/domain"
	websocket "esp32/src/internal/websocket/application"

	"github.com/gin-gonic/gin"
)

type CreateStatusFoodController struct {
	createStatusFood *application.CreateStatusFood
	wsService    *websocket.WebSocketService
    cageRepo     cages.CageRepository
}

func NewCreateStatusFoodController(
	createStatusFood *application.CreateStatusFood, 	
	wsService	 *websocket.WebSocketService,
    cageRepo 	 cages.CageRepository, 
	)*CreateStatusFoodController {
	return &CreateStatusFoodController{
		createStatusFood: createStatusFood, 
		wsService:  wsService,
        cageRepo:    cageRepo,
	}
}

func (h *CreateStatusFoodController) Create(c *gin.Context) {
	var foodRequest domain.Food
	if err := c.ShouldBindJSON(&foodRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Creando estatus de alimento desde HTTP: %+v\n", foodRequest)

	err := h.createStatusFood.Execute(foodRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "estatus de alimento creada correctamente", "food": foodRequest})
}

func (h *CreateStatusFoodController) ProcessFood(food domain.Food) error {
    fmt.Printf("Procesando estatus de alimento desde AMQP: %+v\n", food)
    
    if err := h.createStatusFood.Execute(food); err != nil {
        return err
    }
    
    cage, err := h.cageRepo.GetCageByID(food.IDHamster)
    if err != nil {
        return err
    }
    
    if err := h.wsService.NotifyUser(cage.Idusuario, gin.H{
        "event": "new_food",
        "data":  food,
        "cage_id": food.IDHamster,
        "timestamp": time.Now().Unix(),
    }); err != nil {
        log.Printf("Error notificando al usuario %d: %v", cage.Idusuario, err)
    }
    
    return nil
}