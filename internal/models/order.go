package models

type Order struct {
	ID                int
	Total             float64
	Status            string
	ExternalReference string // pending, approved, rejected
	CreatedAt         string
}
