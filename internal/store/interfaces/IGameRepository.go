package interfaces

import "github.com/rdsalakhov/game-keys-store/internal/model"

type IGameRepository interface {
	Create(*model.Game) error
	Find(int) (*model.Game, error)
}
