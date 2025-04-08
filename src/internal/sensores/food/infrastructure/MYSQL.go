package infrastructure

import (
	"database/sql"
	amqpConsumer "esp32/src/consumer_amqp"
	"esp32/src/internal/sensores/food/domain"
	"fmt"
)

type FoodRepo struct {
	db           *sql.DB
	amqpConsumer *amqpConsumer.RabbitMQConsumer
}

func NewFoodRepo(db *sql.DB, amqpConsumer *amqpConsumer.RabbitMQConsumer) *FoodRepo {
	return &FoodRepo{
		db:           db,
		amqpConsumer: amqpConsumer,
	}
}

func (r *FoodRepo) CreateStatusFood(food domain.Food) error {
	query := "INSERT INTO alimento (idhamster, alimento, porcentaje, hora_registro) VALUES (?, ?, ?, NOW())"
	_, err := r.db.Exec(query, food.IDHamster, food.Alimento, food.Porcentaje)
	if err != nil {
		return fmt.Errorf("error al guardar el estatus de alimento: %w", err)
	}

	return nil
}

func (r *FoodRepo) GetByHamster(IDHamster string) ([]domain.Food, error) {
	query := "SELECT idalimento, idhamster, alimento, porcentaje, hora_registro FROM alimento WHERE idhamster = ?"
	rows, err := r.db.Query(query, IDHamster)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estatus: %w", err)
	}
	defer rows.Close()

	var foods []domain.Food
	for rows.Next() {
		var food domain.Food
		if err := rows.Scan(&food.IDalimento, &food.IDHamster, &food.Alimento, &food.Porcentaje, &food.HoraRegistro); err != nil {
			return nil, fmt.Errorf("error al escanear estatus de alimento: %w", err)
		}
		foods = append(foods, food)
	}

	return foods, nil
}
