package commands

import (
	"fmt"
	"main/exchanges/mexc"
	"os"
	"strings"
)

type ExchangeClient interface {
	CheckConnection()
	GetBalanceUSD() float64
	GetLastPriceBTC() float64
	SetBaseURL(url string)
	CreateOrder(side, price, quantity string) ([]byte, error)
	GetOrderById(id string) ([]byte, error)
	IsFilled(id string) bool
	CancelOrder(orderID string) ([]byte, error)
}

func GetClientByExchange(exchangeArg ...string) ExchangeClient {

	var ex string
	if len(exchangeArg) > 0 {
		ex = exchangeArg[0]
	} else {
		ex = os.Getenv("EXCHANGE")
	}
	ex = strings.ToUpper(ex)

	var client ExchangeClient

	switch ex {
	case "MEXC":
		client = mexc.NewClient()
		client.SetBaseURL("https://api.mexc.co")
	default:
		fmt.Println("Unsupported exchange:", ex)
		os.Exit(0)
	}

	return client
}
