package middleware

import (
	auth "esp32/src/core"
	tokenService "esp32/src/internal/services/auth/domain"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenService tokenService.TokenService, authRepo *auth.AuthRepository, requiredType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token de autorización requerido"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "formato de token inválido"})
			return
		}

		userID, userType, err := tokenService.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
			return
		}

		if requiredType != "" && userType != requiredType {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "acceso no autorizado"})
			return
		}

		user, err := authRepo.FindUserByID(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error al verificar usuario"})
			return
		}

		c.Set("userID", userID)
		c.Set("userType", userType)
		c.Set("user", user)
		c.Next()
	}
}
