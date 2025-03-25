package application

import (
	"esp32/src/internal/motion/domain"
)
type GetByHamster struct {
	repo domain.MotionRepository
}

func NewGetByHamster(repo domain.MotionRepository) *GetByHamster {
	return &GetByHamster{repo: repo}
}

func (cp *GetByHamster) Execute(IDHamster int32) ([]domain.Motion, error){
	return cp.repo.GetByHamster(IDHamster)	
}