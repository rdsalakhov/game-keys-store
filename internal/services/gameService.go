package services

import (
	"github.com/rdsalakhov/game-keys-store/internal/model"
	"github.com/rdsalakhov/game-keys-store/internal/store/interfaces"
)

type GameService struct {
	Store interfaces.IStore
}

func (service *GameService) AddGame(game *model.Game) error {
	if err := service.Store.Game().Create(game); err != nil {
		return err
	}
	return nil
}
