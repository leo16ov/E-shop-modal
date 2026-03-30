package handlers

import "e-shop-modal/internal/server"

func JSONError(ctx *server.Context, code int, message string) {
	ctx.JSONResponse(code, map[string]string{
		"message": message,
	})

}
