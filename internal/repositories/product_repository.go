package repositories

import (
	"database/sql"
	"e-shop-modal/internal/models"
	"fmt"
)

type Store interface {
	GetAll() ([]*models.Product, error)
	GetByID(id int) (*models.Product, error)
	Create(product *models.Product) (*models.Product, error)
	Delete(id int) error
	Update(id int, product *models.Product) (*models.Product, error)
}

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (s *ProductRepository) GetAll() ([]*models.Product, error) {
	q := `SELECT id_producto, tipo, precio, cantidad, talles, colores FROM Producto`
	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		product := models.Product{}
		if err := rows.Scan(&product.ID, &product.Tipo, &product.Precio, &product.Cantidad,
			&product.Talles, &product.Colores); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	return products, nil
}

func (s *ProductRepository) GetByID(id int) (*models.Product, error) {
	q := `SELECT id_producto, tipo, precio, cantidad, talles, colores FROM Producto WHERE id_producto = $1`
	product := models.Product{}
	err := s.db.QueryRow(q, id).Scan(&product.ID, &product.Tipo, &product.Precio, &product.Cantidad, &product.Talles, &product.Colores)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *ProductRepository) Create(product *models.Product) (*models.Product, error) {
	q := `INSERT INTO Producto(tipo, precio, cantidad, talles, colores) VALUES($1, $2, $3, $4, $5) RETURNING id_producto`
	err := s.db.QueryRow(q, product.Tipo, product.Precio, product.Cantidad, product.Talles, product.Colores).Scan(&product.ID)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductRepository) Update(id int, product *models.Product) (*models.Product, error) {
	q := `UPDATE producto SET tipo= $1, precio= $2, cantidad= $3, talles= $4, colores= $5 WHERE id_producto = $6`
	result, err := s.db.Exec(q, product.Tipo, product.Precio, product.Cantidad, product.Talles, product.Colores, id)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("Producto no encontrado")
	}
	product.ID = id
	return product, nil
}

func (s *ProductRepository) Delete(id int) error {
	q := `DELETE FROM Producto WHERE id_producto = $1`
	result, err := s.db.Exec(q, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("producto no encontrado")
	}
	return nil
}
