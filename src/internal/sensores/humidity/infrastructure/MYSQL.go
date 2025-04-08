package infrastructure

import (
	"database/sql"
	amqpConsumer "esp32/src/consumer_amqp"
	"esp32/src/internal/sensores/humidity/domain"
	"fmt"
)

type HumidityRepo struct {
	db           *sql.DB
	amqpConsumer *amqpConsumer.RabbitMQConsumer
}

func NewHumidityRepo(db *sql.DB, amqpConsumer *amqpConsumer.RabbitMQConsumer) *HumidityRepo {
	return &HumidityRepo{
		db:           db,
		amqpConsumer: amqpConsumer,
	}
}

func (r *HumidityRepo) CreateHumidity(humidity domain.Humidity) error {
	query := "INSERT INTO humedad (idhamster, humedad, hora_registro) VALUES (?, ?, NOW())"
	_, err := r.db.Exec(query, humidity.IDHamster, humidity.Humedad)
	if err != nil {
		return fmt.Errorf("error al guardar humedad: %w", err)
	}

	return nil
}

func (r *HumidityRepo) GetByHamster(IDHamster string) ([]domain.Humidity, error) {
	query := "SELECT idhumedad, idhamster, humedad, hora_registro FROM humedad WHERE idhamster = ?"
	rows, err := r.db.Query(query, IDHamster)
	if err != nil {
		return nil, fmt.Errorf("error al obtener humedades: %w", err)
	}
	defer rows.Close()

	var humidities []domain.Humidity
	for rows.Next() {
		var humidity domain.Humidity
		if err := rows.Scan(&humidity.IDhumedad, &humidity.IDHamster, &humidity.Humedad, &humidity.HoraRegistro); err != nil {
			return nil, fmt.Errorf("error al escanear humedad: %w", err)
		}
		humidities = append(humidities, humidity)
	}

	return humidities, nil
}
