package model

type Seller struct {
	ID                int    `json:"id"`
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"-"`
	URL               string `json:"url"`
	Account           string `json:"account"`
}
