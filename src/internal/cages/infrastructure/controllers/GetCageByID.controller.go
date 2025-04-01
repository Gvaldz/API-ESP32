package controllers

import (
	"esp32/src/internal/cages/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetCageByIDController struct {
	getCageByID *application.GetCageByID
}

func NewGetCageByIDController(getCageByID *application.GetCageByID) *GetCageByIDController {
	return &GetCageByIDController{getCageByID: getCageByID}
}

func (h *GetCageByIDController) GetCageByID(c *gin.Context) {
    iduser := c.Param("id")
    idInt, err := strconv.Atoi(iduser)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID de jaula inválido"})
        return
    }

    user, err := h.getCageByID.Execute(int32(idInt))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "data": user,
    })
}