package commands

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"main/exchanges/mexc"
	"math"
	"os"
	"strconv"
)

type ExchangeClient interface {
	CheckConnection()
	GetBalanceUSDT() float64
	GetLastPriceBTC() float64
	SetBaseURL(url string)
	CreateOrder(side string, price string, quantity string) (string, error)
}

func CalcAmountUSDT(freeBalance float64, percentStr string) float64 {
	percent, err := strconv.ParseFloat(percentStr, 64)
	if err != nil {
		log.Fatal(err)
	}
	return percent * freeBalance / 100
}

func CalcAmountBTC(availableUSDT, priceBTC float64) float64 {
	return availableUSDT / priceBTC
}

func FormatSmallFloat(quantity float64) string {
	return fmt.Sprintf("%.6f", quantity)
}

var client ExchangeClient

func New() {
	exchange := os.Getenv("EXCHANGE")
	percent := os.Getenv("PERCENT")

	buyOffset, _ := strconv.ParseFloat(os.Getenv("BUY_OFFSET"), 64)
	buyOffset = math.Abs(buyOffset)

	sellOffset, _ := strconv.ParseFloat(os.Getenv("SELL_OFFSET"), 64)
	sellOffset = math.Abs(sellOffset)

	switch exchange {
	case "MEXC":
		client = mexc.NewClient()
		client.SetBaseURL("https://api.mexc.co")
	default:
		fmt.Println("Unsupported exchange:", exchange)
		os.Exit(0)
	}

	client.CheckConnection()

	freeBalance := client.GetBalanceUSDT()
	fmt.Println("Free USDT Balance:", freeBalance)
	if freeBalance < 10 {
		color.Red("At least 10$ needed")
		os.Exit(0)
	}

	btcPrice := client.GetLastPriceBTC()
	fmt.Println("BTC Price:", btcPrice)
	fmt.Println("")

	newCycleUSDT := CalcAmountUSDT(freeBalance, percent)
	fmt.Println("USDT for this new cycle:", newCycleUSDT)

	newCycleBTC := CalcAmountBTC(newCycleUSDT, btcPrice)
	fmt.Println("BTC for this new cycle:", newCycleBTC)

	newCycleBTCFormated := FormatSmallFloat(newCycleBTC)
	fmt.Println("BTCFormated for this new cycle:", newCycleBTCFormated)

	buyPrice := btcPrice - buyOffset
	fmt.Println("buy price:", buyPrice)

	sellPrice := btcPrice + sellOffset
	fmt.Println("sell price:", sellPrice)

	// Prepare Order
	buyPriceStr := fmt.Sprintf("%.2f", buyPrice)

	orderID, err := client.CreateOrder("BUY", buyPriceStr, newCycleBTCFormated)
	if err != nil {
		fmt.Println("Order failed:", err)
	} else {
		fmt.Println("Order placed successfully, ID:", orderID)
	}

}
