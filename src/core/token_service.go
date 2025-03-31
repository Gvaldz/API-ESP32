package core

import (
    "errors"
    "os"
    "time"
    "esp32/src/internal/auth/domain"
    
    "github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
    secretKey []byte
}

func NewJWTService() *JWTService {
    key := []byte(os.Getenv("JWT_SECRET"))
    if len(key) == 0 {
        panic("JWT_SECRET not set")
    }
    return &JWTService{secretKey: key}
}

func (s *JWTService) GenerateToken(idusuario int32, correo string) (domain.Token, error) {
    expiresAt := time.Now().Add(24 * time.Hour).Unix()
    
    claims := jwt.MapClaims{
        "idusarrio": idusuario,
        "correo":   correo,
        "exp":     expiresAt,
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(s.secretKey)
    if err != nil {
        return domain.Token{}, err
    }
    
    return domain.Token{
        Token:     tokenString,
        ExpiresAt: expiresAt,
    }, nil
}

func (s *JWTService) ValidateToken(tokenString string) (int32, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return s.secretKey, nil
    })
    
    if err != nil {
        return 0, err
    }
    
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        idusuario := int32(claims["user_id"].(float64))
        return idusuario, nil
    }
    
    return 0, errors.New("invalid token")
}