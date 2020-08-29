package services

import (
	"github.com/rdsalakhov/game-keys-store/internal/model"
	"github.com/rdsalakhov/game-keys-store/internal/store/interfaces"
	"time"
)

type PaymentService struct {
	Store interfaces.IStore
}

func (service *PaymentService) CreateSession(keyID int, name string, email string, address string) (int, error) {
	session := &model.PaymentSession{
		KeyID:           keyID,
		Date:            time.Now(),
		CustomerName:    name,
		CustomerEmail:   email,
		CustomerAddress: address,
	}
	if err := service.Store.PaymentSession().Create(session); err != nil {
		return 0, err
	}
	return session.ID, nil
}

func (service *PaymentService) DeleteSession(id int) error {
	err := service.Store.PaymentSession().DeleteByID(id)
	return err
}
