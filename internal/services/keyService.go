package services

import (
	"github.com/rdsalakhov/game-keys-store/internal/model"
	"github.com/rdsalakhov/game-keys-store/internal/store/interfaces"
)

type KeyService struct {
	Store interfaces.IStore
}

func (service *KeyService) AddKeysToGame(gameID int, keys *[]model.Key) error {
	for _, key := range *keys {
		key.GameID = gameID
		key.Status = model.KeyStatusEnum(0)
		if err := service.Store.Key().Create(&key); err != nil {
			return err
		}
	}
	return nil
}
