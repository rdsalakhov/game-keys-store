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
		key.Status = model.KeyStatusAvailable
		if err := service.Store.Key().Create(&key); err != nil {
			return err
		}
	}
	return nil
}

func (service *KeyService) FindAvailableKey(gameID int) (*model.Key, error) {
	key, err := service.Store.Key().FindSingleAvailable(gameID)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (service *KeyService) FindByGameID(gameID int) (*[]model.Key, error) {
	keys, err := service.Store.Key().FindByGameID(gameID)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (service *KeyService) MarkOnHold(keyID int) error {
	err := service.Store.Key().UpdateStatus(keyID, model.KeyStatusOnHold)
	return err
}
