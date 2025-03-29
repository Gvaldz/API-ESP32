package application

import (
	"esp32/src/internal/breeds/domain"
)

type GetBreedByID struct {
	repo domain.BreedRepository
}

func NewGetBreedByID(repo domain.BreedRepository) *GetBreedByID {
	return &GetBreedByID{repo: repo}
}

func (cp *GetBreedByID) Execute(IDRaza int32) (domain.Breed, error){
	return cp.repo.GetBreedByID(IDRaza)	
}	