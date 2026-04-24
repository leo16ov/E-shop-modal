package repositories

import (
	"database/sql"
	"e-shop-modal/internal/models"
	"e-shop-modal/internal/server"
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

func (r *UserRepository) Create(c *server.Context, user *models.User) error {
	q := `INSERT INTO Usuario(nombre, apellido, email, contrasena, telefono, dni, provider, rol) 
		VALUES($1, $2, $3, $4, $5, $6, 'app', 'Admin') RETURNING id_usuario`
	err := r.db.QueryRowContext(c.Context(), q, user.Nombre, user.Apellido, user.Email, user.Contrasena, user.Telefono, user.DNI).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("Error al crear usuario %w", err)
	}
	return nil
}

func (r *UserRepository) CreateUserOAuth(c *server.Context, user *models.User) error {
	q := `INSERT INTO Usuario(nombre, apellido, email, provider, rol)
          VALUES($1, $2, $3, $4, 'Cliente') RETURNING id_usuario, rol`
	err := r.db.QueryRowContext(c.Context(), q, user.Nombre, user.Apellido, user.Email, user.Provider).Scan(&user.ID, &user.Rol)
	if err != nil {
		return fmt.Errorf("Error al crear usuario %w", err)
	}
	return nil
}

func (r *UserRepository) EmailExists(c *server.Context, email string) (bool, error) {
	var count int
	q := "SELECT COUNT(id_usuario) FROM Usuario WHERE email= $1"

	err := r.db.QueryRowContext(c.Context(), q, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("Error al verificar email: %w", err)
	}
	return count > 0, nil
}

func (r *UserRepository) GetByEmail(c *server.Context, email string) (*models.User, error) {
	var user models.User
	var contrasena sql.NullString

	q := "SELECT id_usuario, nombre, apellido, contrasena, rol, email, provider FROM Usuario WHERE email= $1"

	err := r.db.QueryRowContext(c.Context(), q, email).Scan(&user.ID, &user.Nombre, &user.Apellido, &contrasena, &user.Rol, &user.Email, &user.Provider)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener el usuario: %w", err)
	}

	user.Contrasena = contrasena.String // si es NULL, queda como ""
	return &user, nil
}
