package domain

type UserRepository interface {
	CreateUser(User) (User, error)
    GetUserByID(IdUsuario int32) (User, error)
	UpdateUser(IdUsuario int32, user User) error
	UpdatePassword(IdUsuario int32, password string) error
}
