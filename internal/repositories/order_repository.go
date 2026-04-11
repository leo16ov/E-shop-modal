package repositories

import (
	"database/sql"
	"e-shop-modal/internal/models"
	"fmt"
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
	orden := &models.Order{}
	query := `INSERT INTO orden (total, estado) VALUES ($1, $2) RETURNING id_orden`

	err := r.db.QueryRow(query, total, "pending").Scan(&orden.ID)
	if err != nil {
		return nil, err
	}
	orden.Total = total
	orden.Status = "pending"

	return orden, nil
}

func (r *OrderRepository) SetExternalReference(orderID int, ref string) error {
	query := `UPDATE orden SET referencia_externa = $1 WHERE id_orden = $2`

	result, err := r.db.Exec(query, ref, orderID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("Orden no encontrada")
	}
	return nil
}

func (r *OrderRepository) UpdateStatus(orderID int, status string) error {
	query := `UPDATE orden SET estado = $1 WHERE id_orden = $2`

	result, err := r.db.Exec(query, status, orderID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("Orden no encontrada")
	}

	return nil
}

func (r *OrderRepository) GetByID(id int) (*models.Order, error) {
	query := `SELECT id_orden, total, estado, referencia_externa, fecha_emision
		FROM orden WHERE id_orden = $1`

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
	query := `SELECT id_orden, total, estado, referencia_externa, fecha_emision
		FROM orden WHERE referencia_externa = $1`

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
