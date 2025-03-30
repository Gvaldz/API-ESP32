package application

import (
	"esp32/src/internal/users/domain"
)

type GetUserByID struct {
	repo domain.UserRepository
}

func NewGetUserByID(repo domain.UserRepository) *GetUserByID {
	return &GetUserByID{repo: repo}
}

func (cp *GetUserByID) Execute(IDuser int32) (domain.User, error){
	return cp.repo.GetUserByID(IDuser)	
}	