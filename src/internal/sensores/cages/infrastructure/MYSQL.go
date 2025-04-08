package infrastructure

import (
	"database/sql"
	"esp32/src/internal/sensores/cages/domain"
	"fmt"
)

type CageRepo struct {
	db *sql.DB
}

func NewCageRepo(db *sql.DB) *CageRepo {
	return &CageRepo{
		db: db}
}

func (r *CageRepo) CreateCage(cage domain.Cage) error {
	query := "INSERT INTO jaulas (idjaulas, idusuarios, nombre_hamster) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, cage.Idjaula, cage.Idusuario, cage.Nombre_hamster)
	if err != nil {
		return fmt.Errorf("error al guardar jaula: %w", err)
	}

	return nil
}

func (r *CageRepo) GetAllCages() ([]domain.Cage, error) {
	query := "SELECT idjaulas, idusuarios, nombre_hamster FROM jaulas"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener jaulas: %w", err)
	}
	defer rows.Close()

	var cages []domain.Cage
	for rows.Next() {
		var cage domain.Cage
		if err := rows.Scan(&cage.Idjaula, &cage.Idusuario, &cage.Nombre_hamster); err != nil {
			return nil, fmt.Errorf("error al escanear jaulas: %w", err)
		}
		cages = append(cages, cage)
	}

	return cages, nil
}

func (r *CageRepo) GetCageByID(idjaula string) (domain.Cage, error) {
	var cage domain.Cage
	query := "SELECT idjaulas, idusuarios, nombre_hamster FROM jaulas WHERE idjaulas = ?"
	err := r.db.QueryRow(query, idjaula).Scan(&cage.Idjaula, &cage.Idusuario, &cage.Nombre_hamster)
	if err != nil {
		return cage, fmt.Errorf("error al obtener jaula: %w", err)
	}
	return cage, nil
}

func (r *CageRepo) GetCagesByUser(Idusuario int32) ([]domain.Cage, error) {
	query := "SELECT idjaulas, idusuarios, nombre_hamster FROM jaulas WHERE idusuarios = ?"
	rows, err := r.db.Query(query, Idusuario)
	if err != nil {
		return nil, fmt.Errorf("error al obtener jaulas: %w", err)
	}
	defer rows.Close()

	var cages []domain.Cage
	for rows.Next() {
		var cage domain.Cage
		if err := rows.Scan(&cage.Idjaula, &cage.Idusuario, &cage.Nombre_hamster); err != nil {
			return nil, fmt.Errorf("error al obtener jaula: %w", err)
		}
		cages = append(cages, cage)
	}

	return cages, nil
}

func (r *CageRepo) UpdateCage(id string, cage domain.Cage) error {
	query := "UPDATE jaulas SET idusuarios = ?, nombre_hamster = ? WHERE idjaulas = ?"
	result, err := r.db.Exec(query, cage.Idusuario, cage.Nombre_hamster, id)
	if err != nil {
		return fmt.Errorf("error al actualizar jaula: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar actualización: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("jaula no encontrada")
	}

	return nil
}
