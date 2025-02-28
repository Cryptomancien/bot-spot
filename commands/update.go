package commands

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/fatih/color"
	"log"
	"main/database"
	"main/exchanges/mexc"
	"os"
	"strconv"
)

func Update() {
	exchange := os.Getenv("EXCHANGE")

	switch exchange {
	case "MEXC":
		client = mexc.NewClient()
		client.SetBaseURL("https://api.mexc.co")
	default:
		fmt.Println("Unsupported exchange:", exchange)
		os.Exit(0)
	}

	lastPrice := client.GetLastPriceBTC()
	//fmt.Println("Last price:", lastPrice)

	docs := database.List()
	for _, doc := range docs {

		id := doc.Get("_id")
		idInt := doc.Get("idInt")
		idString := id.(string)

		status := doc.Get("status")
		quantity := doc.Get("quantity")
		//buyPrice := doc.Get("buyPrice")
		buyId := doc.Get("buyId")
		sellPrice := doc.Get("sellPrice")
		sellId := doc.Get("sellId")

		if status == "buy" {
			order, _ := client.GetOrderById((buyId).(string))
			isFilled := client.IsFilled(string(order))
			if !isFilled {
				color.Blue("%d Order Buy still active - %s", idInt, buyId)
			} else {
				color.Green("BUY Order Filled:")

				// Check sell price > last price
				sellPrice := (sellPrice).(float64)

				if lastPrice > sellPrice {
					newSellPrice := sellPrice + 100
					fmt.Println("New sell price: ", newSellPrice)
					newSellPriceStr := strconv.FormatFloat(newSellPrice, 'f', 2, 64)

					database.FindCycleByIdAndUpdate(idString, "sellPrice", newSellPriceStr)
					fmt.Println("New sell price updated: ")
				}

				// Place sell order
				quantity := (quantity).(float64)
				quantityStr := strconv.FormatFloat(quantity, 'f', 6, 64)

				doc := database.GetById(idString)
				sellPrice = doc.Get("sellPrice").(float64)
				sellPriceStr := strconv.FormatFloat(sellPrice, 'f', 6, 64)

				bytes, err := client.CreateOrder("SELL", sellPriceStr, quantityStr)
				if err != nil {
					log.Fatal(err)
					return
				}
				orderId, _, _, err := jsonparser.Get(bytes, "orderId")

				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("New sell Order:", string(bytes))

				database.FindCycleByIdAndUpdate(idString, "status", "sell")
				database.FindCycleByIdAndUpdate(idString, "sellId", string(orderId))
			}
		} else if status == "sell" {
			order, _ := client.GetOrderById((sellId).(string))
			isFilled := client.IsFilled(string(order))

			if !isFilled {
				fmt.Printf("%s %s %s\n",
					color.WhiteString("%d", idInt),
					color.CyanString("Order Sell still active -"),
					color.YellowString("%s", sellId),
				)
			} else {
				database.FindCycleByIdAndUpdate(idString, "status", "completed")
				fmt.Printf("%s%s\n",
					color.WhiteString("%d", idInt),
					color.GreenString(" Cycle successfully completed"),
				)
			}
		}
	}
}
