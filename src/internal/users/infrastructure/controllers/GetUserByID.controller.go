package controllers

import (
	"esp32/src/internal/users/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetByUserIDController struct {
	getByUserID *application.GetUserByID
}

func NewGetByUserIDController(getByUserID *application.GetUserByID) *GetByUserIDController {
	return &GetByUserIDController{getByUserID: getByUserID}
}

func (h *GetByUserIDController) GetByUserID(c *gin.Context) {
	iduser := c.Param("iduser")
	idInt, err := strconv.Atoi(iduser)
	user, err := h.getByUserID.Execute(int32(idInt))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
