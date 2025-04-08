package domain

type MotionRepository interface {
    GetByHamster(IDHamster string) ([]Motion, error)
	CreateMotion(Motion) error
}
