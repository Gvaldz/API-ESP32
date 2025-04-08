package application

import (
	"esp32/src/internal/sensores/cages/domain"
)

type UpdateCage struct {
	cageRepo domain.CageRepository
}

func NewUpdateCage(cageRepo domain.CageRepository) *UpdateCage {
	return &UpdateCage{cageRepo: cageRepo}
}

func (uc *UpdateCage) Execute(id string, cage domain.Cage) error {
	return uc.cageRepo.UpdateCage(id, cage)
}
