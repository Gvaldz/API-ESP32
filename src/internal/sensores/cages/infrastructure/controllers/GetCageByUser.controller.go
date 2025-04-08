package controllers

import (
	"esp32/src/internal/sensores/cages/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetCagesByUserController struct {
	getByUser *application.GetCagesByUser
}

func NewGetCagesByUserController(getByUser *application.GetCagesByUser) *GetCagesByUserController {
	return &GetCagesByUserController{getByUser: getByUser}
}

func (h *GetCagesByUserController) GetByUser(c *gin.Context) {
	iduser := c.Param("id")
	idInt, err := strconv.Atoi(iduser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	user, err := h.getByUser.Execute(int32(idInt))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}
