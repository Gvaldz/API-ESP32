package domain

type TemperatureRepository interface {
    GetByHamster(IDHamster int32) ([]Temperature, error)
	CreateTemperature(Temperature) error
}
