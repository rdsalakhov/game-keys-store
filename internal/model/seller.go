package model

import "golang.org/x/crypto/bcrypt"

type Seller struct {
	ID                int    `json:"id"`
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"-"`
	URL               string `json:"url"`
	Account           string `json:"account"`
}

func (seller *Seller) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(seller.EncryptedPassword), []byte(password)) == nil
}

func (seller *Seller) BeforeCreate() error {
	if len(seller.Password) > 0 {
		enc, err := encryptString(seller.Password)
		if err != nil {
			return err
		}

		seller.EncryptedPassword = enc
	}
	return nil
}

func (seller *Seller) HidePassword() {
	seller.Password = ""
}

func encryptString(s string) (string, error) {
	enc, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(enc), nil
}
