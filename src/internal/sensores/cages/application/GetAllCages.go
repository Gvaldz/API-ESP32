package application

import (
	"esp32/src/internal/sensores/cages/domain"
)

type GetAllCages struct {
	repo domain.CageRepository
}

func NewGetAllCages(repo domain.CageRepository) *GetAllCages {
	return &GetAllCages{repo: repo}
}

func (cp *GetAllCages) Execute() ([]domain.Cage, error) {
	return cp.repo.GetAllCages()
}
