package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	"esp32/src/internal/services/fcm"
	"esp32/src/internal/users/application"
	"esp32/src/internal/users/infrastructure/controllers"
	middleware "esp32/src/server/middleware"
)

type UserDependencies struct {
	DB           *sql.DB
	AMQP         *core.AMQPConnection
	Hasher       *core.BcryptHasher
	UserRepo     *core.UserRepository
	FCMSender    *fcm.FCMSender
	AuthRepo     *core.AuthRepository
	TokenService *core.JWTService
}

func NewUserDependencies(db *sql.DB, amqp *core.AMQPConnection, hasher *core.BcryptHasher, fcmSender *fcm.FCMSender, tokenService *core.JWTService, authRepo *core.AuthRepository, userRepo *core.UserRepository) *UserDependencies {

	return &UserDependencies{
		DB:           db,
		AMQP:         amqp,
		Hasher:       hasher,
		FCMSender:    fcmSender,
		TokenService: tokenService,
		AuthRepo:     authRepo,
		UserRepo:     userRepo,
	}
}

func (d *UserDependencies) GetRoutes() *UserRoutes {
	createUserUseCase := application.NewCreateUser(d.UserRepo, d.Hasher)
	getAllUserUseCase := application.NewGetAllUsers(d.UserRepo)
	getUserUseCase := application.NewGetUserByID(d.UserRepo)
	updateUserUseCase := application.NewUpdateUser(d.UserRepo)
	updatePasswordUseCase := application.NewUpdatePassword(d.UserRepo, d.Hasher)
	deleteUserUseCase := application.NewDeleteUser(d.UserRepo)
	fcmController := controllers.NewFCMController(d.UserRepo)

	createUserController := controllers.NewCreateUserController(createUserUseCase)
	getUsersController := controllers.NewGetAllUsersController(getAllUserUseCase)
	getUserController := controllers.NewGetByUserIDController(getUserUseCase)
	updateUserController := controllers.NewUpdateUserController(updateUserUseCase)
	updatePasswordController := controllers.NewUpdatePasswordController(updatePasswordUseCase)
	deleteUserController := controllers.NewDeleteUserController(deleteUserUseCase)
	authMiddleware := middleware.AuthMiddleware(d.TokenService, d.AuthRepo, "usuario")

	return NewUserRoutes(
		createUserController,
		getUsersController,
		getUserController,
		updateUserController,
		updatePasswordController,
		deleteUserController,
		fcmController,
		authMiddleware,
	)
}
