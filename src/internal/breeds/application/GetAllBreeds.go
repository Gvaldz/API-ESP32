package application

import (
	"esp32/src/internal/breeds/domain"
)

type GetAllBreeds struct {
	repo domain.BreedRepository
}

func NewGetAllBreeds(repo domain.BreedRepository) *GetAllBreeds {
	return &GetAllBreeds{repo: repo}
}

func (cp *GetAllBreeds) Execute() ([]domain.Breed, error){
	return cp.repo.GetAllBreeds()	
}	