package services

import (
	"e-shop-modal/internal/config"
	"e-shop-modal/internal/dto"
	"e-shop-modal/internal/models"
	"e-shop-modal/internal/repositories"
	"e-shop-modal/internal/server"
	"e-shop-modal/internal/utils"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository *repositories.UserRepository
	config     *config.Config
	jwtManager *utils.JWTManager
}

func NewUserService(r *repositories.UserRepository, cfg *config.Config, jwt *utils.JWTManager) *UserService {
	return &UserService{
		repository: r,
		config:     cfg,
		jwtManager: jwt,
	}
}

func (s *UserService) SignUp(c *server.Context, user *models.User) (*models.User, error) {
	exists, err := s.repository.EmailExists(c, user.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("El mail ya esta registrado.")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Contrasena), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Error al hashear la contraseña: %w", err)
	}
	user.Contrasena = string(hashedPassword)
	err = s.repository.Create(c, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) LogIn(c *server.Context, email, contrasena string) (*dto.LoginResponse, error) {
	user, err := s.repository.GetByEmail(c, email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Contrasena), []byte(contrasena))
	if err != nil {
		return nil, fmt.Errorf("Credenciales invalidas")
	}
	token, err := s.jwtManager.GenerateJWT(uint(user.ID), user.Email, user.Rol)
	if err != nil {
		return nil, fmt.Errorf("Error al generar token")
	}

	return &dto.LoginResponse{
		User: &dto.UserLogin{
			ID:       user.ID,
			Nombre:   user.Nombre,
			Apellido: user.Apellido,
			Email:    user.Email,
			Rol:      user.Rol,
		},
		Token: token}, nil
}
