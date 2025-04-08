package core

import (
	"database/sql"
	"esp32/src/internal/services/auth/domain"
	users "esp32/src/internal/users/domain"
	"fmt"
)

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(DB *sql.DB) domain.AuthRepository {
	return &AuthRepository{DB: DB}
}

func (r *AuthRepository) FindUserByEmail(email string) (users.User, error) {
	var user users.User
	query := `SELECT idusuarios, correo, contrasena, tipo FROM usuarios WHERE correo = ?`
	err := r.DB.QueryRow(query, email).Scan(&user.IdUsuario, &user.Correo, &user.Contrasena, &user.Tipo)
	return user, err
}

func (r *AuthRepository) UpdateLastLogin(userID int32) error {
	query := `UPDATE usuarios SET ultimo_login = NOW() WHERE idusuarios = ?`
	_, err := r.DB.Exec(query, userID)
	return err
}

func (r *AuthRepository) FindUserByID(userID int32) (users.User, error) {
	var user users.User
	query := `SELECT idusuarios, correo, contrasena, FCMtoken tipo FROM usuarios WHERE idusuarios = ?`
	err := r.DB.QueryRow(query, userID).Scan(&user.IdUsuario, &user.Correo, &user.Contrasena, &user.Tipo, &user.FCMToken)
	return user, err
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(DB *sql.DB) users.UserRepository {
	return &UserRepository{DB: DB}
}

func (r *UserRepository) CreateUser(user users.User) (users.User, error) {
	result, err := r.DB.Exec(
		"INSERT INTO usuarios (nombre, correo, contrasena) VALUES (?, ?, ?)",
		user.Nombre, user.Correo, user.Contrasena,
	)
	if err != nil {
		return users.User{}, fmt.Errorf("error al crear usuario: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return users.User{}, fmt.Errorf("error al obtener ID: %w", err)
	}

	return users.User{
		IdUsuario: int32(id),
		Nombre:    user.Nombre,
		Correo:    user.Correo,
	}, nil
}

func (r *UserRepository) GetAllUsers() ([]users.User, error) {
	query := "SELECT idusuarios, nombre, correo FROM usuarios"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuarios: %w", err)
	}
	defer rows.Close()

	var usersList []users.User
	for rows.Next() {
		var user users.User
		if err := rows.Scan(&user.IdUsuario, &user.Nombre, &user.Correo); err != nil {
			return nil, fmt.Errorf("error al escanear user: %w", err)
		}
		usersList = append(usersList, user)
	}

	return usersList, nil
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

func (r *UserRepository) GetUserByEmail(email string) (users.User, error) {
	var user users.User
	query := "SELECT idusuarios, nombre, correo, contrasena FROM usuarios WHERE correo = ?"
	err := r.DB.QueryRow(query, email).Scan(
		&user.IdUsuario,
		&user.Nombre,
		&user.Correo,
		&user.Contrasena,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("usuario no encontrado")
		}
		return user, fmt.Errorf("error al obtener usuario por email: %w", err)
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(id int32, user users.User) error {
	query := "UPDATE usuarios SET nombre = ?, correo = ? WHERE idusuarios = ?"
	result, err := r.DB.Exec(query, user.Nombre, user.Correo, id)
	if err != nil {
		return fmt.Errorf("error al actualizar usuario: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar actualización: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("usuario no encontrado")
	}

	return nil
}

func (r *UserRepository) UpdatePassword(id int32, newHashedPassword string) error {
	query := "UPDATE usuarios SET contrasena = ? WHERE idusuarios = ?"
	result, err := r.DB.Exec(query, newHashedPassword, id)
	if err != nil {
		return fmt.Errorf("error al actualizar contraseña: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar actualización de contraseña: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("usuario no encontrado")
	}

	return nil
}

func (r *UserRepository) DeleteUser(id int32) error {
	query := "DELETE FROM usuarios WHERE idusuarios = ?"
	result, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar usuario: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar eliminación: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("usuario no encontrado")
	}

	return nil
}

func (r *UserRepository) UpdateFCMToken(userID string, token string) error {
	query := "UPDATE usuarios SET fcm_token = ? WHERE idusuarios = ?"
	_, err := r.DB.Exec(query, token, userID)
	return err
}

func (r *UserRepository) GetFCMToken(userID int32) (string, error) {
	var token string
	query := "SELECT fcm_token FROM usuarios WHERE idusuarios = ?"
	err := r.DB.QueryRow(query, userID).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}
