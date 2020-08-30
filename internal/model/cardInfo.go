package model

type CardInfo struct {
	Number  string `json:"number"`
	Name    string `json:"name"`
	ExpDate string `json:"exp_date"`
	CVV     string `json:"cvv"`
}
