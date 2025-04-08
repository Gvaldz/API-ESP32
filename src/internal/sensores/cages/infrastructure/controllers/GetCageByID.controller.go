package controllers

import (
	"esp32/src/internal/sensores/cages/application"
	"net/http"

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

	user, err := h.getCageByID.Execute(string(iduser))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}
