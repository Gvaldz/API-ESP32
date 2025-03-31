package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	"esp32/src/internal/users/application"
	"esp32/src/internal/users/infrastructure/controllers"
)

type UserDependencies struct {
	DB       *sql.DB
	AMQP     *core.AMQPConnection
	Hasher   *core.BcryptHasher
	UserRepo *UsersRepo
}

func NewUserDependencies(db *sql.DB, amqp *core.AMQPConnection, hasher *core.BcryptHasher) *UserDependencies {
	userRepo := NewUsersRepo(db) 
	
	return &UserDependencies{
		DB:       db,
		AMQP:     amqp,
		Hasher:   hasher,
		UserRepo: userRepo, 
	}
}

func (d *UserDependencies) GetRoutes() *UserRoutes {
	createUserUseCase := application.NewCreateUser(d.UserRepo, d.Hasher)
	getAllUserUseCase := application.NewGetAllUsers(d.UserRepo)
	getUserUseCase := application.NewGetUserByID(d.UserRepo)
	updateUserUseCase := application.NewUpdateUser(d.UserRepo)
	updatePasswordUseCase := application.NewUpdatePassword(d.UserRepo, d.Hasher)
	deleteUserUseCase := application.NewDeleteUser(d.UserRepo)

	createUserController := controllers.NewCreateUserController(createUserUseCase)
	getUsersController := controllers.NewGetAllUsersController(getAllUserUseCase)
	getUserController := controllers.NewGetByUserIDController(getUserUseCase)
	updateUserController := controllers.NewUpdateUserController(updateUserUseCase)
	updatePasswordController := controllers.NewUpdatePasswordController(updatePasswordUseCase)
	deleteUserController := controllers.NewDeleteUserController(deleteUserUseCase)

	return NewUserRoutes(
		createUserController,
		getUsersController, 
		getUserController, 
		updateUserController, 
		updatePasswordController, 
		deleteUserController,
	)
}