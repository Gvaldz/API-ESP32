package domain

type MotionRepository interface {
    GetByHamster(IDHamster int32) ([]Motion, error)
	CreateMotion(Motion) error
}
