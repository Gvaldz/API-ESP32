package domain

type FoodRepository interface {
	CreateStatusFood(Food) error
    GetByHamster(IDHamster string) ([]Food, error)
}
