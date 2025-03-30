package application

import (
	"esp32/src/internal/users/domain"
	"esp32/src/core"
)

type UpdatePassword struct {
	userRepo domain.UserRepository
	hasher   core.PasswordHasher
}

func NewUpdatePassword(userRepo domain.UserRepository, hasher core.PasswordHasher) *UpdatePassword {
	return &UpdatePassword{
		userRepo: userRepo,
		hasher: hasher,
	}
}

func (uc *UpdatePassword) Execute(id int32, newPassword string) error {
	hashedPassword, err := uc.hasher.Hash(newPassword)
	if err != nil {
		return err
	}
	
	user := domain.User{
		IdUsuario:  id,
		Contrasena: hashedPassword,
	}
	
	return uc.userRepo.UpdatePassword(id, user.Contrasena)
}