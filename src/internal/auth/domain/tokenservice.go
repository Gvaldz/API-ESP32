package domain

type TokenService interface {
    GenerateToken(userID int32, email string) (Token, error)
    ValidateToken(tokenString string) (int32, error)
}
