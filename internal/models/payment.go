package models

import "time"

type PaymentInfo struct {
	ID                int
	Status            string
	TransactionAmount float64
	ExternalReference string
	PaymentMethod     string
	DateApproved      time.Time
}
