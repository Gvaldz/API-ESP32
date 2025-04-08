package controllers

import (
	"esp32/src/internal/sensores/cages/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetAllCagesController struct {
	getAllCages *application.GetAllCages
}

func NewGetAllCagesController(getAllCages *application.GetAllCages) *GetAllCagesController {
	return &GetAllCagesController{getAllCages: getAllCages}
}

func (h *GetAllCagesController) GetAllCages(c *gin.Context) {
	cages, err := h.getAllCages.Execute()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cages)
}
