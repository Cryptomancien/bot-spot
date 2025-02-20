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
	GetBalanceUSDT() float32
	GetLastPriceBTC() float32
	SetBaseURL(url string)
}

func CalcAmountUSDT(freeBalance float32, percentStr string) float32 {
	percent64, err := strconv.ParseFloat(percentStr, 64)
	if err != nil {
		log.Fatal(err)
	}

	percent := float32(percent64)

	return percent * freeBalance / 100
}

func CalcAmountBTC(availableUSDT float32, priceBTC float32) float32 {
	return availableUSDT / priceBTC
}

func FormatSmallFloat(quantity float32) string {
	return fmt.Sprintf("%.6f", quantity)
}

var client ExchangeClient

func New() {
	var exchange = os.Getenv("EXCHANGE")
	var percent = os.Getenv("PERCENT")

	var buyOffset, _ = strconv.ParseFloat(os.Getenv("BUY_OFFSET"), 32)

	buyOffset = math.Abs(buyOffset)

	var sellOffset, _ = strconv.ParseFloat(os.Getenv("SELL_OFFSET"), 32)

	sellOffset = math.Abs(sellOffset)

	switch exchange {
	case "MEXC":
		client = mexc.NewClient()
		client.SetBaseURL("https://api.mexc.co")
	default:
		fmt.Println("Unsupported exchange:", os.Getenv("EXCHANGE"))
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
	fmt.Println("USDT for this new cycle: ", newCycleUSDT)

	newCycleBTC := CalcAmountBTC(newCycleUSDT, btcPrice)
	fmt.Println("BTC for this new cycle: ", newCycleBTC)

	newCycleBTCFormated := FormatSmallFloat(newCycleBTC)
	fmt.Println("BTCFormated for this new cycle: ", newCycleBTCFormated)

}
