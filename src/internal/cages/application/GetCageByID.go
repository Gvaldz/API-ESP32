package application

import (
	"esp32/src/internal/cages/domain"
)

type GetCageByID struct {
	repo domain.CageRepository
}

func NewGetCageByID(repo domain.CageRepository) *GetCageByID {
	return &GetCageByID{repo: repo}
}

func (cp *GetCageByID) Execute(IDRaza int32) (domain.Cage, error){
	return cp.repo.GetCageByID(IDRaza)	
}	