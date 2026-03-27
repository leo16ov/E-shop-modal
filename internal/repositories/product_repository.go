package repositories

import (
	"database/sql"
	"e-shop-modal/internal/models"
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
	q := `SELECT id_producto, tipo, precio, cantidad, talles, colores FROM Producto WHERE id_producto = ?`
	product := models.Product{}
	err := s.db.QueryRow(q, id).Scan(&product.ID, &product.Tipo, &product.Precio, &product.Cantidad, &product.Talles, &product.Colores)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *ProductRepository) Create(product *models.Product) (*models.Product, error) {
	q := `INSERT INTO Producto(tipo, precio, cantidad, talles, colores) VALUES(?, ?, ?, ?, ?)`
	resp, err := s.db.Exec(q, product.Tipo, product.Precio, product.Cantidad, product.Talles, product.Colores)
	if err != nil {
		return nil, err
	}

	id, err := resp.LastInsertId()
	if err != nil {
		return nil, err
	}

	product.ID = int(id)
	return product, nil
}

func (s *ProductRepository) Update(id int, product *models.Product) (*models.Product, error) {
	q := `UPDATE producto SET tipo= ?, precio= ?, cantidad= ?, talles= ?, colores= ? WHERE id_producto = ?`
	_, err := s.db.Exec(q, product.Tipo, product.Precio, product.Cantidad, product.Talles, product.Colores, id)
	if err != nil {
		return nil, err
	}
	product.ID = id
	return product, nil
}

func (s *ProductRepository) Delete(id int) error {
	q := `DELETE FROM Producto WHERE id_producto = ?`
	_, err := s.db.Exec(q, id)
	if err != nil {
		return err
	}
	return nil
}
