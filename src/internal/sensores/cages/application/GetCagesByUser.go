package application

import (
	"esp32/src/internal/sensores/cages/domain"
)

type GetCagesByUser struct {
	repo domain.CageRepository
}

func NewGetCagesByUser(repo domain.CageRepository) *GetCagesByUser {
	return &GetCagesByUser{repo: repo}
}

func (cp *GetCagesByUser) Execute(IDHamster int32) ([]domain.Cage, error) {
	return cp.repo.GetCagesByUser(IDHamster)
}
