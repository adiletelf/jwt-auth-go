package middleware

import (
	"net/http"

	"github.com/adiletelf/jwt-auth-go/internal/config"
	"github.com/adiletelf/jwt-auth-go/internal/token"
	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	cfg, _ := config.New()
	return func(c *gin.Context) {
		err := token.RequestTokenValid(c.Request, "accessToken", cfg.ApiSecret)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
