package infrastructure

import (
	"database/sql"
	"esp32/src/internal/cages/application"
	"esp32/src/internal/cages/infrastructure/controllers"
)

type CageDependencies struct {
	DB     *sql.DB
}

func NewCageDependencies(db *sql.DB) *CageDependencies {
	return &CageDependencies{DB:db}
}

func (d *CageDependencies) GetRoutes() *CageRoutes {
	cageRepo := NewCageRepo(d.DB)

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

	return NewCageRoutes(createCageController,getAllCagesController, getCageByIdController, getCageByUserController,updateCageController)
}