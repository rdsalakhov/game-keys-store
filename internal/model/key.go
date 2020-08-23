package model

type KeyStatusEnum int

const (
	available KeyStatusEnum = iota
	on_hold
	sold
)

type Key struct {
	ID        int           `json:"id"`
	KeyString string        `json:"key_string"`
	GameID    int           `json:"game_id"`
	SellerID  int           `json:"seller_id"`
	Status    KeyStatusEnum `json:"status"`
}
