package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Game struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	SellerID    int     `json:"seller_id"`
	Price       float32 `json:"price"`
	OnSale      bool    `json:"on_sale"`
}

func (game *Game) Validate() error {
	return validation.ValidateStruct(
		game,
		validation.Field(&game.Title, validation.Required),
		validation.Field(&game.Price, validation.Required, validation.Min(0.01)),
		validation.Field(&game.SellerID, validation.Required),
	)
}
