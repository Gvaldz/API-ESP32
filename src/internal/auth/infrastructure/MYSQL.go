package infrastructure

import (
    "database/sql"
    user "esp32/src/internal/users/domain"
)

type AuthRepositoryImpl struct {
    db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepositoryImpl {
    return &AuthRepositoryImpl{db: db}
}

func (r *AuthRepositoryImpl) FindUserByEmail(email string) (user.User, error) {
    var user user.User
    query := `SELECT idusuarios, correo, contrasena, tipo FROM usuarios WHERE correo = ?`
    err := r.db.QueryRow(query, email).Scan(&user.IdUsuario, &user.Correo, &user.Contrasena, &user.Tipo)
    return user, err
}

func (r *AuthRepositoryImpl) UpdateLastLogin(userID int32) error {
    query := `UPDATE usuarios SET ultimo_login = NOW() WHERE idusuarios = ?`
    _, err := r.db.Exec(query, userID)
    return err
}

func (r *AuthRepositoryImpl) FindUserByID(userID int32) (user.User, error) {
    var user user.User
    query := `SELECT idusuarios, correo, contrasena, tipo FROM usuarios WHERE idusuarios = ?`
    err := r.db.QueryRow(query, userID).Scan(&user.IdUsuario, &user.Correo, &user.Contrasena, &user.Tipo)
    return user, err
}