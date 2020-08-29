package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"time"
)

type PaymentSession struct {
	ID              int       `json:"id"`
	Price           float64   `json:"price"`
	Date            time.Time `json:"date"`
	KeyID           int       `json:"key_id"`
	CustomerName    string    `json:"customer_name"`
	CustomerEmail   string    `json:"customer_email"`
	CustomerAddress string    `json:"customer_address"`
	IsPerformed     bool      `json:"is_performed"`
}

func (session *PaymentSession) Validate() error {
	return validation.ValidateStruct(
		session,
		validation.Field(&session.CustomerName, validation.Required),
		validation.Field(&session.CustomerEmail, validation.Required, is.Email),
		validation.Field(&session.CustomerAddress, validation.Required),
	)
}
