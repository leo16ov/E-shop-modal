package repositories

import (
	"database/sql"
	"e-shop-modal/internal/models"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *models.User) error {
	q := `INSERT INTO Usuario(nombre, apellido, email, contrasena, telefono, dni, rol) 
		VALUES($1, $2, $3, $4, $5, $6, 'Admin') RETURNING id_usuario`
	err := r.db.QueryRow(q, user.Nombre, user.Apellido, user.Email, user.Contrasena, user.Telefono, user.DNI).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("Error al crear usuario %w", err)
	}
	return nil
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	var count int
	q := "SELECT COUNT(id_usuario) FROM Usuario WHERE email= $1"

	err := r.db.QueryRow(q, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("Error al verificar email: %w", err)
	}
	return count > 0, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	q := "SELECT id_usuario, nombre, apellido, contrasena, rol, email FROM Usuario WHERE email= $1"

	err := r.db.QueryRow(q, email).Scan(&user.ID, &user.Nombre, &user.Apellido, &user.Contrasena, &user.Rol, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener el usuario")
	}
	return &user, nil
}
