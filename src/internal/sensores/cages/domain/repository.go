package domain

type CageRepository interface{
	CreateCage(Cage) error
	GetAllCages() ([]Cage, error)
	GetCageByID(idcage string) (Cage, error)
	GetCagesByUser(iduser int32)([]Cage, error)
	UpdateCage(idcage string, cage Cage) error
}