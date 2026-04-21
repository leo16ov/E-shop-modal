package handlers

import (
	"e-shop-modal/internal/config"
	"e-shop-modal/internal/server"
	"e-shop-modal/internal/services"
	"e-shop-modal/internal/utils"
	"net/http"

	"golang.org/x/oauth2"
)

type OAuthHandler struct {
	service           *services.OAuthService
	config            *config.Config
	jwtManager        *utils.JWTManager
	googleOAuthConfig *oauth2.Config
}

func NewOAuthHandler(s *services.OAuthService, cfg *config.Config, jwt *utils.JWTManager, oauthConfig *oauth2.Config) *OAuthHandler {
	return &OAuthHandler{
		service:           s,
		config:            cfg,
		jwtManager:        jwt,
		googleOAuthConfig: oauthConfig,
	}
}

// Redirige a Google
func (h *OAuthHandler) GoogleLogin(c *server.Context) {
	state, err := utils.GenerateRandomState()
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "Error al generar state")
		return
	}
	// Guarda el state en una cookie httpOnly
	http.SetCookie(c.RWriter, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		MaxAge:   600, // 10 minutos, suficiente para completar el flujo
		HttpOnly: true,
		Secure:   h.config.Debug != "dev", // true en producción
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	url := h.googleOAuthConfig.AuthCodeURL(state)
	http.Redirect(c.RWriter, c.Request, url, http.StatusTemporaryRedirect) //Corregir esto
}

// Callback
func (h *OAuthHandler) GoogleCallback(c *server.Context) {
	// Lee state de la cookie
	cookie, err := c.Request.Cookie("oauth_state")
	if err != nil {
		JSONError(c, http.StatusBadRequest, "State cookie no encontrada")
		return
	}

	// Compara con el state que mandó Google
	stateRecibido := c.Request.URL.Query().Get("state")
	if cookie.Value != stateRecibido {
		JSONError(c, http.StatusUnauthorized, "State inválido, posible ataque CSRF")
		return
	}

	// Invalida la cookie inmediatamente (uso único)
	http.SetCookie(c.RWriter, &http.Cookie{
		Name:   "oauth_state",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})

	// Continua con el flujo normal
	code := c.Request.URL.Query().Get("code")
	user, err := h.service.LoginWithGoogle(c, h.googleOAuthConfig, h.config.OAuthIDClient, h.config.OAuthSecretClient, h.config.NotificationURL, code)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "Error en login")
		return
	}

	token, err := h.jwtManager.GenerateJWT(uint(user.ID), user.Email, user.Rol)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "Error generando token")
		return
	}

	c.JSONResponse(http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}
