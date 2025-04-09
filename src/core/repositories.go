package core

import (
	"database/sql"
	users "esp32/src/internal/users/domain"
	"fmt"
)


type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(DB *sql.DB) users.UserRepository {
	return &UserRepository{DB: DB}
}


func (r *UserRepository) GetUserByID(iduser int32) (users.User, error) {
	if r.DB == nil {
		return users.User{}, fmt.Errorf("database connection is nil")
	}

	var user users.User
	query := "SELECT idusuarios, nombre, correo, FCMtoken FROM usuarios WHERE idusuarios = ?"
	err := r.DB.QueryRow(query, iduser).Scan(&user.IdUsuario, &user.Nombre, &user.Correo, &user.FCMToken)
	if err != nil {
		return user, fmt.Errorf("error al obtener usuario: %w", err)
	}
	return user, nil
}
