package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	"esp32/src/internal/users/application"
	"esp32/src/internal/users/infrastructure/controllers"
)

type UserDependencies struct {
	DB     *sql.DB
	AMQP   *core.AMQPConnection
	Hasher *core.BcryptHasher
}

func NewUserDependencies(db *sql.DB, amqp *core.AMQPConnection, hasher *core.BcryptHasher) *UserDependencies {
	return &UserDependencies{
		DB:     db,
		AMQP:   amqp,
		Hasher: hasher,
	}
}

func (d *UserDependencies) GetRoutes() *UserRoutes {
	userRepo := NewUsersRepo(d.DB)

	createUserUseCase := application.NewCreateUser(userRepo, d.Hasher)
	getAllUserUseCase := application.NewGetAllUsers(userRepo)
	getUserUseCase := application.NewGetUserByID(userRepo)
	updateUserUseCase := application.NewUpdateUser(userRepo)
	updatePasswordUseCase := application.NewUpdatePassword(userRepo, d.Hasher)
	deleteUserUseCase := application.NewDeleteUser(userRepo)

	createUserController := controllers.NewCreateUserController(createUserUseCase)
	getUsersController := controllers.NewGetAllUsersController(getAllUserUseCase)
	getUserController := controllers.NewGetByUserIDController(getUserUseCase)
	updateUserController := controllers.NewUpdateUserController(updateUserUseCase)
	updatePasswordController := controllers.NewUpdatePasswordController(updatePasswordUseCase)
	deleteUserController := controllers.NewDeleteUserController(deleteUserUseCase)

	return NewUserRoutes(createUserController,getUsersController, getUserController, updateUserController, updatePasswordController, deleteUserController)
}