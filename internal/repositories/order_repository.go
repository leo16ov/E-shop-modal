package repositories

import (
	"database/sql"
	"e-shop-modal/internal/models"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(total float64) (*models.Order, error) {
	query := `INSERT INTO orders (total, status) VALUES (?, ?)`

	result, err := r.db.Exec(query, total, "pending")
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	order := models.Order{
		ID:     int(id),
		Total:  total,
		Status: "pending",
	}

	return &order, nil
}
func (r *OrderRepository) SetExternalReference(orderID int, ref string) error {
	query := `UPDATE orders SET external_reference = ? WHERE id = ?`

	_, err := r.db.Exec(query, ref, orderID)
	return err
}

func (r *OrderRepository) UpdateStatus(orderID int, status string) error {
	query := `UPDATE orders SET status = ? WHERE id = ?`

	_, err := r.db.Exec(query, status, orderID)
	return err
}

func (r *OrderRepository) GetByID(id int) (*models.Order, error) {
	query := `SELECT id, total, status, external_reference, created_at
		FROM orders WHERE id = ?`

	var order models.Order

	err := r.db.QueryRow(query, id).Scan(
		&order.ID, &order.Total, &order.Status,
		&order.ExternalReference,
		&order.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetByExternalReference(ref string) (*models.Order, error) {
	query := `SELECT id, total, status, external_reference, created_at
		FROM orders WHERE external_reference = ?`

	var order models.Order

	err := r.db.QueryRow(query, ref).Scan(
		&order.ID, &order.Total, &order.Status,
		&order.ExternalReference,
		&order.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
