package infrastructure

import (
    "database/sql"
    "esp32/src/core"
    "esp32/src/internal/auth/application"
    "esp32/src/internal/auth/infrastructure/controllers"
    userRepo "esp32/src/internal/users/infrastructure"
)

type AuthDependencies struct {
    DB       *sql.DB
    Hasher   *core.BcryptHasher
    UserRepo *userRepo.UsersRepo
}

func NewAuthDependencies(db *sql.DB, hasher *core.BcryptHasher, userRepo *userRepo.UsersRepo) *AuthDependencies {
    return &AuthDependencies{
        DB:       db,
        Hasher:   hasher,
        UserRepo: userRepo,
    }
}

func (d *AuthDependencies) GetRoutes() *AuthRoutes {
    authRepo := NewAuthRepository(d.DB)
    tokenService := core.NewJWTService()
    
    loginUC := application.NewLogin(
        authRepo,
        d.UserRepo,
        tokenService,
        d.Hasher,
    )
    
    loginController := controllers.NewLoginController(loginUC)
    
    return NewAuthRoutes(loginController)
}