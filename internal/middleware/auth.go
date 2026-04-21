package middleware

import (
	"e-shop-modal/internal/server"
	"e-shop-modal/internal/utils"
	"errors"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	jwtManager *utils.JWTManager
}

func NewAuthMiddleware(jwtManager *utils.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{jwtManager: jwtManager}
}

func (m *AuthMiddleware) Authentication(next func(*server.Context)) func(*server.Context) {
	return func(c *server.Context) {
		authorization := c.Request.Header.Get("Authorization")

		tokenStr, err := tokenFromAuthentication(authorization)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.Send(err.Error())
			return // ← también faltaba este return
		}

		claims, err := m.jwtManager.ValidateJWT(tokenStr)
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

func tokenFromAuthentication(authorization string) (string, error) {
	if authorization == "" {
		return "", errors.New("authorization header requerido")
	}
	parts := strings.SplitN(authorization, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("formato inválido, usar: Bearer <token>")
	}
	return parts[1], nil
}
