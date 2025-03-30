package infrastructure

import (
	"database/sql"
	"esp32/src/internal/users/domain"
	"fmt"
)

type UsersRepo struct {
	db *sql.DB
}

func NewUsersRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

func (r *UsersRepo) CreateUser(user domain.User) (domain.User, error) {
    result, err := r.db.Exec(
        "INSERT INTO usuarios (nombre, correo, contrasena) VALUES (?, ?, ?)",
        user.Nombre, user.Correo, user.Contrasena,
    )
    if err != nil {
        return domain.User{}, fmt.Errorf("error al crear usuario: %w", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return domain.User{}, fmt.Errorf("error al obtener ID: %w", err)
    }

    return domain.User{
        IdUsuario: int32(id),
        Nombre:    user.Nombre,
        Correo:    user.Correo,
    }, nil
}

func (r *UsersRepo) GetUserByID(iduser int32) (domain.User, error) {
	var user domain.User
	query := "SELECT id_usuario, nombre, correo FROM usuarios WHERE id_usuario = ?"
	err := r.db.QueryRow(query, iduser).Scan(&user.IdUsuario, &user.Nombre, &user.Correo)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("usuario no encontrado")
		}
		return user, fmt.Errorf("error al obtener usuario: %w", err)
	}
	return user, nil
}

func (r *UsersRepo) GetUserByEmail(email string) (domain.User, error) {
	var user domain.User
	query := "SELECT id_usuario, nombre, correo, contrasena FROM usuarios WHERE correo = ?"
	err := r.db.QueryRow(query, email).Scan(
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

func (r *UsersRepo) UpdateUser(id int32, user domain.User) error {
	query := "UPDATE usuarios SET nombre = ?, correo = ? WHERE id_usuario = ?"
	result, err := r.db.Exec(query, user.Nombre, user.Correo, id)
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

func (r *UsersRepo) UpdatePassword(id int32, newHashedPassword string) error {
	query := "UPDATE usuarios SET contrasena = ? WHERE id_usuario = ?"
	result, err := r.db.Exec(query, newHashedPassword, id)
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

func (r *UsersRepo) DeleteUser(id int32) error {
	query := "DELETE FROM usuarios WHERE id_usuario = ?"
	result, err := r.db.Exec(query, id)
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