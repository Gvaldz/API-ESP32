package infrastructure

import (
	"database/sql"
	amqpConsumer "esp32/src/consumer_amqp"
	"esp32/src/internal/sensores/motion/domain"
	"fmt"
)

type MotionRepo struct {
	db           *sql.DB
	amqpConsumer *amqpConsumer.RabbitMQConsumer
}

func NewMotionRepo(db *sql.DB, amqpConsumer *amqpConsumer.RabbitMQConsumer) *MotionRepo {
	return &MotionRepo{
		db:           db,
		amqpConsumer: amqpConsumer,
	}
}

func (r *MotionRepo) CreateMotion(motion domain.Motion) error {
	query := "INSERT INTO movimiento (idhamster, movimiento, hora_registro) VALUES (?, ?, NOW())"
	_, err := r.db.Exec(query, motion.IDHamster, motion.Movimiento)
	if err != nil {
		return fmt.Errorf("error al guardar movimiento: %w", err)
	}

	return nil
}

func (r *MotionRepo) GetByHamster(IDHamster string) ([]domain.Motion, error) {
	query := "SELECT idmovimiento, idhamster, movimiento, hora_registro FROM movimiento WHERE idhamster = ?"
	rows, err := r.db.Query(query, IDHamster)
	if err != nil {
		return nil, fmt.Errorf("error al obtener movimientos: %w", err)
	}
	defer rows.Close()

	var motions []domain.Motion
	for rows.Next() {
		var motion domain.Motion
		if err := rows.Scan(&motion.IDMovimiento, &motion.IDHamster, &motion.Movimiento, &motion.HoraRegistro); err != nil {
			return nil, fmt.Errorf("error al escanear movimiento: %w", err)
		}
		motions = append(motions, motion)
	}

	return motions, nil
}
