package domain

type BreedRepository interface{
	GetAllBreeds() ([]Breed, error)
	GetBreedByID(IDRaza int32) (Breed, error)
}