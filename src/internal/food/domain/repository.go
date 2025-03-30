package domain

type FoodRepository interface {
	CreateStatusFood(Food) error
    GetByHamster(IDHamster int32) ([]Food, error)
}
