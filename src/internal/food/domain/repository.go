package domain

type FoodRepository interface {
    GetByHamster(IDHamster int32) ([]Food, error)
	CreateStatusFood(Food) error
}
