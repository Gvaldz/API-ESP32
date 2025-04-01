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
    iduser := c.Param("id")
    idInt, err := strconv.Atoi(iduser)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
        return
    }

    user, err := h.getByUserID.Execute(int32(idInt))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "data": user,
    })
}