package services

import (
	"e-shop-modal/internal/dto"
	"e-shop-modal/internal/models"
	"e-shop-modal/internal/repositories"
	"e-shop-modal/internal/utils"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository *repositories.UserRepository
}

func NewUserService(r *repositories.UserRepository) *UserService {
	return &UserService{
		repository: r,
	}
}

func (s *UserService) SignUp(user *models.User) (*models.User, error) {
	exists, err := s.repository.EmailExists(user.Email)
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
	err = s.repository.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) LogIn(email, contrasena string) (*dto.LoginResponse, error) {
	user, err := s.repository.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Contrasena), []byte(contrasena))
	if err != nil {
		return nil, fmt.Errorf("Credenciales invalidas")
	}
	token, err := utils.GenerateJWT(uint(user.ID), user.Email, user.Rol)
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
