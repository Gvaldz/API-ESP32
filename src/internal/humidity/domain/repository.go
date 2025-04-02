package domain

type HumidityRepository interface {
    GetByHamster(IDHamster string) ([]Humidity, error)
	CreateHumidity(Humidity) error
}
