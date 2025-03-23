package infrastructure

import (
	"database/sql"
	"esp32/src/internal/temperatura/domain"
	"fmt"
)

type TemperatureRepo struct {
	db           *sql.DB
	amqpConsumer *AMQPConsumer // AMQPConsumer para recibir mensajes
}

// Constructor
func NewTemperatureRepo(db *sql.DB, amqpConsumer *AMQPConsumer) *TemperatureRepo {
	return &TemperatureRepo{
		db:           db,
		amqpConsumer: amqpConsumer,
	}
}

// Método para guardar temperatura en la base de datos
func (r *TemperatureRepo) CreateTemperature(temperature domain.Temperature) error {
	query := "INSERT INTO temperatura (idhamster, temperatura, hora_registro) VALUES (?, ?, NOW())"
	_, err := r.db.Exec(query, temperature.IDHamster, temperature.Temperatura)
	if err != nil {
		return fmt.Errorf("error al guardar temperatura: %w", err)
	}

	return nil
}

// Método para obtener las temperaturas de un hámster por su ID
func (r *TemperatureRepo) GetByHamster(IDHamster int32) ([]domain.Temperature, error) {
	query := "SELECT idtemperatura, idhamster, temperatura, hora_registro FROM temperatura WHERE idhamster = ?"
	rows, err := r.db.Query(query, IDHamster)
	if err != nil {
		return nil, fmt.Errorf("error al obtener temperaturas: %w", err)
	}
	defer rows.Close()

	var temperatures []domain.Temperature
	for rows.Next() {
		var temperature domain.Temperature
		if err := rows.Scan(&temperature.IDtemperatura, &temperature.IDHamster, &temperature.Temperatura, &temperature.HoraRegistro); err != nil {
			return nil, fmt.Errorf("error al escanear temperatura: %w", err)
		}
		temperatures = append(temperatures, temperature)
	}

	return temperatures, nil
}
