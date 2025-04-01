package controllers

import (
	"esp32/src/internal/users/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeleteUserController struct {
	deleteUser *application.DeleteUser
}

func NewDeleteUserController(delete *application.DeleteUser) *DeleteUserController {
	return &DeleteUserController{
		deleteUser: delete,
	}
}

func (h *DeleteUserController) Delete(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = h.deleteUser.Execute(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado correctamente"})
}