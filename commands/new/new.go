package new

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/fatih/color"
	"log"
	"main/commands"
	"main/database"
	"main/exchanges/mexc"
	"math"
	"os"
	"strconv"
)

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

var client commands.ExchangeClient

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
	color.White("Free USD Balance: %.2f", freeBalance)
	if freeBalance < 10 {
		color.Red("At least 10$ needed")
		os.Exit(0)
	}

	btcPrice := client.GetLastPriceBTC()

	fmt.Printf("%s %s\n",
		color.CyanString("BTC Price"),
		color.YellowString("%.2f", btcPrice),
	)

	newCycleUSDC := CalcAmountUSD(freeBalance, percent)

	fmt.Printf("%s %s\n",
		color.CyanString("USD for this new cycle:"),
		color.YellowString("%.2f", newCycleUSDC),
	)

	newCycleBTC := CalcAmountBTC(newCycleUSDC, btcPrice)
	newCycleBTCFormated := FormatSmallFloat(newCycleBTC)
	fmt.Printf("%s %s\n",
		color.CyanString("BTC for this new cycle:"),
		color.YellowString(newCycleBTCFormated),
	)

	buyPrice := btcPrice - buyOffset
	fmt.Printf("%s %s\n",
		color.CyanString("Buy Price"),
		color.YellowString("%.2f", buyPrice),
	)

	sellPrice := btcPrice + sellOffset
	fmt.Printf("%s %s\n",
		color.CyanString("Sell Price"),
		color.YellowString("%.2f", sellPrice),
	)

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
