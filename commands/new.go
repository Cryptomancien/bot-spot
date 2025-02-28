package commands

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/fatih/color"
	"log"
	"main/database"
	"main/exchanges/mexc"
	"math"
	"os"
	"strconv"
)

type ExchangeClient interface {
	CheckConnection()
	GetBalanceUSD() float64
	GetLastPriceBTC() float64
	SetBaseURL(url string)
	CreateOrder(side, price, quantity string) ([]byte, error)
	GetOrderById(id string) ([]byte, error)
	IsFilled(id string) bool
}

func CalcAmountUSD(freeBalance float64, percentStr string) float64 {
	percent, err := strconv.ParseFloat(percentStr, 64)
	if err != nil {
		log.Fatal(err)
	}
	return percent * freeBalance / 100
}

func CalcAmountBTC(availableUSD, priceBTC float64) float64 {
	return availableUSD / priceBTC
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

	freeBalance := client.GetBalanceUSD()
	fmt.Println("Free USDC Balance:", freeBalance)
	if freeBalance < 10 {
		color.Red("At least 10$ needed")
		os.Exit(0)
	}

	btcPrice := client.GetLastPriceBTC()
	fmt.Println("BTC Price:", btcPrice)
	fmt.Println("")

	newCycleUSDC := CalcAmountUSD(freeBalance, percent)
	fmt.Println("USDC for this new cycle:", newCycleUSDC)

	newCycleBTC := CalcAmountBTC(newCycleUSDC, btcPrice)
	fmt.Println("BTC for this new cycle:", newCycleBTC)

	newCycleBTCFormated := FormatSmallFloat(newCycleBTC)
	fmt.Println("BTCFormated for this new cycle:", newCycleBTCFormated)

	buyPrice := btcPrice - buyOffset
	fmt.Println("buy price:", buyPrice)

	sellPrice := btcPrice + sellOffset
	fmt.Println("sell price:", sellPrice)

	// Prepare Order
	buyPriceStr := fmt.Sprintf("%.2f", buyPrice)

	body, err := client.CreateOrder("BUY", buyPriceStr, newCycleBTCFormated)
	if err != nil {
		color.Red("Order failed:", err)
		os.Exit(0)
	}

	orderId, _, _, err := jsonparser.Get(body, "orderId")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("OrderId:", string(orderId))

	// Insert in database
	cycle := database.Cycle{
		Exchange:  "mexc", // Todo dynamic
		Status:    "buy",
		Quantity:  newCycleBTC,
		BuyPrice:  buyPrice,
		BuyId:     string(orderId),
		SellPrice: sellPrice,
		SellId:    "",
	}
	database.NewCycle(&cycle)

	color.Green("New Cycle successfully inserted in database:")
}
