package domain

type UserRepository interface {
	CreateUser(User) (User, error)
	GetAllUsers() ([]User, error)
    GetUserByID(IdUsuario int32) (User, error)
	UpdateUser(IdUsuario int32, user User) error
	UpdatePassword(IdUsuario int32, password string) error
	DeleteUser(IdUsuario int32) error
    UpdateFCMToken(userID string, token string) error
    GetFCMToken(userID int32) (string, error)
}
