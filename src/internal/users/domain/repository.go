package domain

type UserRepository interface {
    GetUserByID(IdUsuario int32) (User, error)
}
