package infrastructure

import (
	"database/sql"
	"esp32/src/internal/breeds/domain"
	"fmt"
)

type BreadRepo struct {
	db           *sql.DB
}

func NewBreadRepo(db *sql.DB) *BreadRepo {
	return &BreadRepo{
		db:           db	}
}

func (r *BreadRepo) GetByID(IDRaza int32) (domain.Breed, error) {
	var breed domain.Breed
	query := "SELECT idrazas, raza WHERE idrazas = ?"
	err := r.db.QueryRow(query, IDRaza).Scan(&breed.IDRaza, &breed.Raza)
	if err != nil {
		return breed, fmt.Errorf("error al obtener raza: %w", err)
	}
	return breed, nil
}

func (r *BreadRepo) GetAll() ([]domain.Breed, error) { 
	query := "SELECT idrazas, raza FROM razas"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener razas: %w", err)
	}
	defer rows.Close()

	var breeds []domain.Breed
	for rows.Next() {
		var breed domain.Breed
		if err := rows.Scan(&breed.IDRaza, &breed.Raza); err != nil {
			return nil, fmt.Errorf("error al escanear razas: %w", err)
		}
		breeds = append(breeds, breed)
	}

	return breeds, nil
}