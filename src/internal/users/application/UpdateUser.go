package application

import (
	"esp32/src/internal/users/domain"
)

type UpdateUser struct {
	userRepo domain.UserRepository
}

func NewUpdateUser(userRepo domain.UserRepository) *UpdateUser {
	return &UpdateUser{userRepo: userRepo}
}

func (uc *UpdateUser) Execute(id int32, user domain.User) error {
	user.Contrasena = ""
	return uc.userRepo.UpdateUser(id, user)
}