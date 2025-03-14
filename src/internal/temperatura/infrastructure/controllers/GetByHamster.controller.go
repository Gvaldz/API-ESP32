package controllers

import (
	"esp32/src/internal/temperatura/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetByHamsterController struct {
	getByHamster *application.GetByHamster
}

func NewGetByHamsterController(getByHamster *application.GetByHamster) *GetByHamsterController {
	return &GetByHamsterController{getByHamster: getByHamster}
}

func (h *GetByHamsterController) GetByHamster(c *gin.Context) {
	idHamster := c.Param("idHamster")
	idInt, err := strconv.Atoi(idHamster)
	temperatures, err := h.getByHamster.Execute(int32(idInt))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, temperatures)
}
