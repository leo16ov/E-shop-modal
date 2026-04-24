package services

import (
	"e-shop-modal/internal/dto"
	"e-shop-modal/internal/models"
	"e-shop-modal/internal/repositories"
	"e-shop-modal/internal/server"
	"encoding/json"
	"errors"
	"io"

	"golang.org/x/oauth2"
)

type OAuthService struct {
	repository        *repositories.UserRepository
	googleOAuthConfig *oauth2.Config
}

func NewOAuthService(r *repositories.UserRepository, oauthConfig *oauth2.Config) *OAuthService {
	return &OAuthService{
		repository:        r,
		googleOAuthConfig: oauthConfig,
	}
}

func (s *OAuthService) LoginWithGoogle(c *server.Context, googleOAuthConfig *oauth2.Config, code string) (*dto.UserOAuth, error) {
	// Intercambia code por token
	token, err := s.googleOAuthConfig.Exchange(c.Context(), code)
	if err != nil {
		return nil, err
	}

	// Crea cliente con token
	client := s.googleOAuthConfig.Client(c.Context(), token)

	// Obtiene datos del usuario
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var userGoogle dto.UserGoogle
	if err := json.Unmarshal(body, &userGoogle); err != nil {
		return nil, err
	}

	// Valida si exite ese email en la DB
	exist, err := s.repository.EmailExists(c, userGoogle.Email)
	if err != nil {
		return nil, err
	}
	if userGoogle.Email == "" {
		return nil, errors.New("Email no recibido de Google")
	}
	if !exist {
		// Crea usuario si no esta registrado
		newUser := &models.User{
			Email:    userGoogle.Email,
			Nombre:   userGoogle.GivenName,
			Apellido: userGoogle.FamilyName,
			Provider: "google",
		}

		err = s.repository.CreateUserOAuth(c, newUser)
		if err != nil {
			return nil, err
		}
		return &dto.UserOAuth{
			ID:       newUser.ID,
			Nombre:   newUser.Nombre,
			Email:    newUser.Email,
			Rol:      newUser.Rol,
			Provider: newUser.Provider,
		}, nil
	}
	user, err := s.repository.GetByEmail(c, userGoogle.Email)
	if err != nil {
		return nil, err
	}
	return &dto.UserOAuth{
		ID:       user.ID,
		Nombre:   user.Nombre,
		Email:    user.Email,
		Rol:      user.Rol,
		Provider: user.Provider,
	}, nil
}
