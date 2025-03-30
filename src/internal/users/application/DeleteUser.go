package application

import (
	"esp32/src/internal/users/domain"
)

type DeleteUser struct {
	repo domain.UserRepository
}

func NewDeleteUser(repo domain.UserRepository) *DeleteUser {
	return &DeleteUser{repo: repo}
}

func (lp *DeleteUser) Execute(id int32) error{
	return lp.repo.DeleteUser(id)
}