package domain

type CageRepository interface{
	CreateCage(Cage) error
	GetAllCages() ([]Cage, error)
	GetCageByID(idcage int32) (Cage, error)
	GetCagesByUser(iduser int32)([]Cage, error)
	UpdateCage(idcage int32, cage Cage) error
}