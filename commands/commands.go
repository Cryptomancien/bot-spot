package commands

type ExchangeClient interface {
	CheckConnection()
	GetBalanceUSD() float64
	GetLastPriceBTC() float64
	SetBaseURL(url string)
	CreateOrder(side, price, quantity string) ([]byte, error)
	GetOrderById(id string) ([]byte, error)
	IsFilled(id string) bool
}
