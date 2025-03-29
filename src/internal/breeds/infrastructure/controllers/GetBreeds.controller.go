package controllers

import (
	"esp32/src/internal/breeds/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetAllBreedsController struct {
	getAllBreeds *application.GetAllBreeds
}

func NewGetBreedsController(getAllBreeds *application.GetAllBreeds) *GetAllBreedsController {
	return &GetAllBreedsController{getAllBreeds: getAllBreeds}
}

func (h *GetAllBreedsController) GetBreeds(c *gin.Context) {
	breeds, err := h.getAllBreeds.Execute()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, breeds)
}
