package mysqlStore

import (
	"database/sql"
	"github.com/rdsalakhov/game-keys-store/internal/store/interfaces"
)

type Store struct {
	db                       *sql.DB
	gameRepository           interfaces.IGameRepository
	sellerRepository         interfaces.ISellerRepository
	keyRepository            interfaces.IKeyRepository
	paymentSessionRepository interfaces.IPaymentSessionRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (store *Store) Game() interfaces.IGameRepository {
	if store.gameRepository != nil {
		return store.gameRepository
	}

	store.gameRepository = &GameRepository{
		store: store,
	}

	return store.gameRepository
}

func (store *Store) Seller() interfaces.ISellerRepository {
	if store.sellerRepository != nil {
		return store.sellerRepository
	}

	store.sellerRepository = &SellerRepository{
		store: store,
	}

	return store.sellerRepository
}

func (store *Store) Key() interfaces.IKeyRepository {
	if store.keyRepository != nil {
		return store.keyRepository
	}

	store.keyRepository = &KeyRepository{
		store: store,
	}
	return store.keyRepository
}

func (store *Store) PaymentSession() interfaces.IPaymentSessionRepository {
	if store.paymentSessionRepository != nil {
		return store.paymentSessionRepository
	}

	store.paymentSessionRepository = &PaymentSessionRepository{
		store: store,
	}
	return store.paymentSessionRepository
}
