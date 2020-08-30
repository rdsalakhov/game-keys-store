package model

type PaymentInfo struct {
	SellerAccount   string
	SellerURL       string
	GameTitle       string
	Key             string
	TotalAmount     float64
	SellerAmount    float64
	PlatformAmount  float64
	CustomerName    string
	CustomerEmail   string
	CustomerAddress string
}
