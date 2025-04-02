package domain

type TemperatureRepository interface {
    GetByHamster(IDHamster string) ([]Temperature, error)
	CreateTemperature(Temperature) error
}
