package application

import (
    "esp32/src/internal/users/domain"
    "esp32/src/core"
)

type CreateUser struct {
    userRepo domain.UserRepository
    hasher   core.PasswordHasher
}

func NewCreateUser(userRepo domain.UserRepository, hasher core.PasswordHasher) *CreateUser {
    return &CreateUser{
        userRepo: userRepo,
        hasher: hasher,
    }
}

func (uc *CreateUser) Execute(user domain.User) (domain.User, error) {
    hashedPassword, err := uc.hasher.Hash(user.Contrasena)
    if err != nil {
        return domain.User{}, err
    }
    
    user.Contrasena = hashedPassword
    
    createdUser, err := uc.userRepo.CreateUser(user)
    if err != nil {
        return domain.User{}, err
    }
    
    return createdUser, nil
}