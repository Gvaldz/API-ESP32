package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"esp32/src/internal/users/application"
	"esp32/src/internal/users/domain"
)

type CreateUserController struct {
	createUser *application.CreateUser
}

func NewCreateUserController(createUser *application.CreateUser) *CreateUserController {
	return &CreateUserController{createUser: createUser}
}

func (h *CreateUserController) CreateUser(c *gin.Context) {
	var userRequest domain.User
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de usuario inválidos: " + err.Error()})
		return
	}

	createdUser, err := h.createUser.Execute(userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear usuario: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuario creado correctamente",
		"user": gin.H{
			"id":     createdUser.IdUsuario,
			"nombre": createdUser.Nombre,
			"correo": createdUser.Correo,
		},
	})
}