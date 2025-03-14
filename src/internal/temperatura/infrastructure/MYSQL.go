package infrastructure

import (
	"database/sql"
	"esp32/src/internal/temperatura/domain"
	"fmt"
	"log"
)

type TemperatureRepo struct {
	db      *sql.DB
	mqttProducer *MQTTProducer
}

func NewTemperatureRepo(db *sql.DB, mqttProducer *MQTTProducer) *TemperatureRepo {
	return &TemperatureRepo{db: db, mqttProducer: mqttProducer}
}

func (r *TemperatureRepo) CreateTemperature(temperature domain.Temperature) error {
	query := "INSERT INTO temperatura (idhamster, temperatura) VALUES (?, ?)"
	_, err := r.db.Exec(query, temperature.IDHamster, temperature.Temperatura)
	if err != nil {
		return fmt.Errorf("error al guardar temperatura: %w", err)
	}

	if err := r.mqttProducer.SendTemperatureMessage(temperature); err != nil {
		log.Printf("Error enviando mensaje MQTT: %s", err)
	}

	return nil
}

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