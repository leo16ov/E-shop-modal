package models

import "time"

type PaymentInfo struct {
	ID                int       `json:"id"`
	Status            string    `json:"status"`
	TransactionAmount float64   `json:"transaction_amount"`
	ExternalReference string    `json:"external_reference"`
	PaymentMethod     string    `json:"payment_method_id"`
	DateApproved      time.Time `json:"date_approved"`
}
