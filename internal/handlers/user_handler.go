package handlers

import (
	"e-shop-modal/internal/models"
	"e-shop-modal/internal/server"
	"e-shop-modal/internal/services"
	"net/http"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{
		service: s,
	}
}

func (h *UserHandler) HandleSignUp(c *server.Context) {
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSONResponse(http.StatusBadRequest, "JSON invalido")
		return
	}

	if user.Email == "" || user.Contrasena == "" || user.Nombre == "" {
		c.JSONResponse(http.StatusBadRequest, "Email y Contraseña son requeridos")
		return
	}

	created, err := h.service.SignUp(&user)
	if err != nil {
		c.JSONResponse(http.StatusInternalServerError, err.Error())
		return
	}
	created.Contrasena = ""
	c.JSONResponse(http.StatusCreated, created)
}
