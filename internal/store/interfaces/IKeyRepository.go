package interfaces

import "github.com/rdsalakhov/game-keys-store/internal/model"

type IKeyRepository interface {
	Create(seller *model.Key) error
	Find(ID int) (*model.Key, error)
}
