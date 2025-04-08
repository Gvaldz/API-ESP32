package application

import (
	"esp32/src/internal/sensores/cages/domain"
)

type GetCageByID struct {
	repo domain.CageRepository
}

func NewGetCageByID(repo domain.CageRepository) *GetCageByID {
	return &GetCageByID{repo: repo}
}

func (cp *GetCageByID) Execute(IDRaza string) (domain.Cage, error) {
	return cp.repo.GetCageByID(IDRaza)
}
