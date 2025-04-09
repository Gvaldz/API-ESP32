package domain

type TokenService interface {
    ValidateToken(tokenString string) (int32, string, error) 
}
