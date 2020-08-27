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

func (service *GameService) FindByID(id int) (*model.Game, error) {
	game, err := service.Store.Game().Find(id)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (service *GameService) FindAll() ([]*model.Game, error) {
	games, err := service.Store.Game().FindAll()
	if err != nil {
		return nil, err
	}
	return games, nil
}
