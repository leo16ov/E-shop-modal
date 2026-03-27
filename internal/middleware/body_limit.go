package middleware

import (
	"e-shop-modal/internal/server"
	"net/http"
)

func BodyLimitMiddleware(next func(*server.Context)) func(*server.Context) {
	return func(c *server.Context) {
		c.Request.Body = http.MaxBytesReader(c.RWriter, c.Request.Body, 1<<20)
		defer c.Request.Body.Close()

		next(c)
	}
}
