package domain

type HumidityRepository interface {
    GetByHamster(IDHamster int32) ([]Humidity, error)
	CreateHumidity(Humidity) error
}
