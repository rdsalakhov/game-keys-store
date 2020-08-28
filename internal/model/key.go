package model

import validation "github.com/go-ozzo/ozzo-validation"

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
	Status    KeyStatusEnum `json:"status"`
}

func (key *Key) Validate() error {
	return validation.ValidateStruct(
		key,
		validation.Field(&key.KeyString, validation.Required),
	)
}
