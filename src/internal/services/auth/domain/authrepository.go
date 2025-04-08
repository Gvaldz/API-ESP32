package domain

import (
    user "esp32/src/internal/users/domain"
)

type AuthRepository interface {
    FindUserByEmail(email string) (user.User, error)
    UpdateLastLogin(userID int32) error
}