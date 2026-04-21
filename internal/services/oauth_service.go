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

func (s *OAuthService) LoginWithGoogle(c *server.Context, googleOAuthConfig *oauth2.Config, cID, cSecret, URL, code string) (*dto.UserOAuth, error) {
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
	json.Unmarshal(body, &userGoogle)

	// Buscar usuario en DB
	user, err := s.repository.GetByEmail(c, userGoogle.Email)

	if user == nil {
		if userGoogle.Email == "" {
			return nil, errors.New("email no recibido de Google")
		}
		// Crea usuario si no esta registrado
		newUser := &models.User{
			Email:    userGoogle.Email,
			Nombre:   userGoogle.Name,
			Provider: "google",
		}

		err = s.repository.Create(c, newUser)
		if err != nil {
			return nil, err
		}
		return &dto.UserOAuth{
			Nombre: newUser.Nombre,
			Email:  newUser.Email,
		}, nil
	}
	return &dto.UserOAuth{
		Nombre: user.Nombre,
		Email:  user.Email,
	}, nil
}
