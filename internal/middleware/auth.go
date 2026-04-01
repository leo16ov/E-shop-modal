package middleware

import (
	"e-shop-modal/internal/server"
	"e-shop-modal/internal/utils"
	"errors"
	"net/http"
	"strings"
)

func Authentication(next func(*server.Context)) func(*server.Context) {
	return func(c *server.Context) {

		authorization := c.Request.Header.Get("Authorization")

		tokenStr, err := tokenFromAuthentication(authorization)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.Send(err.Error())
		}

		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.Send("Token inválido")
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("rol", claims.Rol)

		next(c)
	}
}

// Valida que el formato sea correcto y devuelve nada mas que el token en string
func tokenFromAuthentication(authorization string) (string, error) {
	if authorization == "" {
		return "", errors.New("Token requerido")
	}

	parts := strings.Split(authorization, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("Formato invalido de token")
	}
	return parts[1], nil
}
