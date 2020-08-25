package interfaces

type IStore interface {
	Game() IGameRepository
	Seller() ISellerRepository
	Key() IKeyRepository
	PaymentSession() IPaymentSessionRepository
}
