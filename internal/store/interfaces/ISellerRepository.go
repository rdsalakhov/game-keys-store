package interfaces

import "github.com/rdsalakhov/game-keys-store/internal/model"

type ISellerRepository interface {
	Create(seller *model.Seller) error
	Find(ID int) (*model.Seller, error)
	FindByEmail(email string) (*model.Seller, error)
}
