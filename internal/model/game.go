package model

type Game struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	OnSale      bool    `json:"on_sale"`
}
