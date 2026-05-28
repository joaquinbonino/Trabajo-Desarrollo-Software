package utils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			Error(c, http.StatusUnauthorized, "token requerido")
			c.Abort()
			return
		}

		claims, err := ValidateToken(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			Error(c, http.StatusUnauthorized, "token inválido o expirado")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("rol", claims.Rol)
		c.Next()
	}
}
