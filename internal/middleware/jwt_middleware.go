package middleware

import (
	"e-shop-modal/internal/server"
	"e-shop-modal/internal/utils"
	"net/http"
	"strings"
)

func JWTMiddleware(next func(*server.Context)) func(*server.Context) {
	return func(c *server.Context) {

		authHeader := c.Request.Header.Get("Authorization")

		if authHeader == "" {
			c.Status(http.StatusUnauthorized)
			c.Send("Token requerido")
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Status(http.StatusUnauthorized)
			c.Send("Formato de token inválido")
			return
		}

		tokenStr := parts[1]

		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.Send("Token inválido")
			return
		}

		// 🔥 SETEO EN TU CONTEXT
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("rol", claims.Rol)

		next(c)
	}
}
