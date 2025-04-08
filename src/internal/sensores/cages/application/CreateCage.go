package application

import (
	"esp32/src/internal/sensores/cages/domain"
)

type CreateCage struct {
	repo domain.CageRepository
}

func NewCreateCage(repo domain.CageRepository) *CreateCage {
	return &CreateCage{repo: repo}
}

func (c *CreateCage) Execute(cage domain.Cage) error {
	return c.repo.CreateCage(cage)
}
