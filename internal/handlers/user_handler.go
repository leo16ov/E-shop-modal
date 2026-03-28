package handlers

import (
	"e-shop-modal/internal/models"
	"e-shop-modal/internal/server"
	"e-shop-modal/internal/services"
	"e-shop-modal/internal/utils"
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

func (h *UserHandler) HandleLogIn(c *server.Context) {
	var req models.User

	err := c.BindJSON(&req)
	if err != nil {
		c.JSONResponse(http.StatusBadRequest, map[string]interface{}{
			"error": "Credenciales requeridas",
		})
		return
	}
	if req.Email == "" || req.Contrasena == "" {
		c.JSONResponse(http.StatusBadRequest, map[string]interface{}{
			"error": "Credenciales inválidas",
		})
		return
	}

	user, err := h.service.LogIn(req.Email, req.Contrasena)
	if err != nil {
		c.JSONResponse(http.StatusUnauthorized, err)
		return
	}
	token, err := utils.GenerateJWT(uint(user.ID), user.Email, user.Rol)
	if err != nil {
		c.JSONResponse(http.StatusInternalServerError, map[string]interface{}{
			"error": "Error generando token",
		})
		return
	}
	user.Contrasena = ""

	c.JSONResponse(http.StatusOK, map[string]interface{}{
		"user":  user,
		"token": token,
	})

}

func (h *UserHandler) Profile(c *server.Context) {
	val := c.Get("user_id")

	userID, ok := val.(uint)
	if !ok {
		c.Status(http.StatusUnauthorized)
		c.Send("No autorizado")
		return
	}

	email, _ := c.Get("email").(string)
	rol, _ := c.Get("rol").(string)

	c.JSONResponse(http.StatusOK, map[string]interface{}{
		"user_id": userID,
		"email":   email,
		"rol":     rol,
	})
}
