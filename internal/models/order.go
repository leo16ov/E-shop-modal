package models

import "time"

type Order struct {
	ID                int
	Total             float64
	Status            string // pending, approved, rejected
	ExternalReference string
	PaymentID         int
	CreatedAt         time.Time
}
