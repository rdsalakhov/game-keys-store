package model

import "time"

type PaymentSession struct {
	ID              int       `json:"id"`
	Price           float64   `json:"price"`
	Date            time.Time `json:"date"`
	KeyID           int       `json:"key_id"`
	CustomerName    string    `json:"customer_name"`
	CustomerEmail   string    `json:"customer_email"`
	CustomerAddress string    `json:"customer_address"`
}
