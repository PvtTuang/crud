package middleware

import (
	"auth/pkg/authentication"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware(authService *authentication.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		valid, err := authService.ValidateToken(tokenStr)
		if err != nil || !valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		ctx.Next()
	}
}
