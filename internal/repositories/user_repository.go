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
		VALUES(?, ?, ?, ?, ?, ?, Admin)`
	resp, err := r.db.Exec(q, user.Nombre, user.Apellido, user.Email, user.Contrasena, user.Telefono, user.DNI)
	if err != nil {
		return fmt.Errorf("Error al crear usuario %w", err)
	}
	id, err := resp.LastInsertId()

	if err != nil {
		return fmt.Errorf("Error al obtener ID del usuario creado %w", err)
	}
	user.ID = int(id)
	return nil
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	var count int
	q := "SELECT COUNT(id_usuario) FROM Usuario WHERE email= ?"

	err := r.db.QueryRow(q, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("Error al verificar email: %w", err)
	}
	return count > 0, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	q := "SELECT id_usuario, nombre, apellido, contrasena, telefono, dni, rol FROM Usuario WHERE email= ?"

	err := r.db.QueryRow(q, email).Scan(&user.ID, &user.Nombre, &user.Apellido, &user.Contrasena, &user.Telefono, &user.DNI, &user.Rol)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener el usuario")
	}
	return &user, nil
}
