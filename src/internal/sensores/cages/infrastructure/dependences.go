package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	"esp32/src/internal/sensores/cages/application"
	"esp32/src/internal/sensores/cages/infrastructure/controllers"
)

type CageDependencies struct {
	DB *sql.DB
}

func NewCageDependencies(db *sql.DB) *CageDependencies {
	return &CageDependencies{DB: db}
}

func (d *CageDependencies) GetRoutes() *CageRoutes {
	cageRepo := NewCageRepo(d.DB)
	tokenService := core.NewJWTService()
	authRepo := core.NewAuthRepository(d.DB).(*core.AuthRepository)

	createCageUseCase := application.NewCreateCage(cageRepo)
	getAllCageUseCase := application.NewGetAllCages(cageRepo)
	getCageUseCase := application.NewGetCageByID(cageRepo)
	getCageByUserUseCase := application.NewGetCagesByUser(cageRepo)
	updateCageUseCase := application.NewUpdateCage(cageRepo)

	createCageController := controllers.NewCreateCageController(createCageUseCase)
	getAllCagesController := controllers.NewGetAllCagesController(getAllCageUseCase)
	getCageByIdController := controllers.NewGetCageByIDController(getCageUseCase)
	getCageByUserController := controllers.NewGetCagesByUserController(getCageByUserUseCase)
	updateCageController := controllers.NewUpdateCageController(updateCageUseCase)

	return NewCageRoutes(createCageController, getAllCagesController, getCageByIdController, getCageByUserController, updateCageController, tokenService, authRepo)
}
