package interfaces

import "github.com/rdsalakhov/game-keys-store/internal/model"

type IKeyRepository interface {
	Create(seller *model.Key) error
	Find(ID int) (*model.Key, error)
	FindSingleAvailable(gameID int) (*model.Key, error)
	UpdateStatus(keyID int, newStatus string) error
}
