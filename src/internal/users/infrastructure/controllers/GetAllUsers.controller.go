package controllers

import (
	"esp32/src/internal/users/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetAllUsersController struct {
	GetAllUsers  *application.GetAllUsers
}

func NewGetAllUsersController(getAllUsers *application.GetAllUsers) *GetAllUsersController {
	return &GetAllUsersController{
		GetAllUsers:  getAllUsers,
	}
}

func (h *GetAllUsersController) GetAll(c *gin.Context) {
	doctors, err := h.GetAllUsers.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, doctors)
}