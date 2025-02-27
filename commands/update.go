package commands

import (
	"fmt"
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
	fmt.Println("Last price:", lastPrice)

	docs := database.List()
	for _, doc := range docs {

		id := doc.Get("_id")
		idString := id.(string)

		status := doc.Get("status")
		quantity := doc.Get("quantity")
		buyPrice := doc.Get("buyPrice")
		buyId := doc.Get("buyId")
		sellPrice := doc.Get("sellPrice")
		sellId := doc.Get("sellId")

		if status == "buy" {
			fmt.Println(id, status, quantity, buyPrice, sellPrice, buyId, sellId)
			order, _ := client.GetOrderById((buyId).(string))
			isFilled := client.IsFilled(string(order))
			if !isFilled {
				fmt.Println("BUY Order Not filled:")
			} else {
				fmt.Println("BUY Order Filled:")

				// Check sell price > last price
				sellPrice := doc.Get("sellPrice").(float64)

				if lastPrice > sellPrice {
					newSellPrice := sellPrice + 100
					fmt.Println("New sell price: ", newSellPrice)
					newSellPriceStr := strconv.FormatFloat(newSellPrice, 'f', 2, 64)

					database.FindCycleByIdAndUpdate(idString, "sellPrice", newSellPriceStr)
					fmt.Println("New sell price updated: ")
				}

				// Place sell order
				quantity := doc.Get("quantity").(float64)
				quantityStr := strconv.FormatFloat(quantity, 'f', 6, 64)
				fmt.Println("New quantity: ", quantityStr)

				doc := database.GetById(idString)
				sellPrice = doc.Get("sellPrice").(float64)
				sellPriceStr := strconv.FormatFloat(sellPrice, 'f', 6, 64)

				log.Print("quantityStr: ", quantityStr)
				return

				bytes, err := client.CreateOrder("SELL", sellPriceStr, quantityStr)
				if err != nil {
					log.Fatal(err)
					return
				}
				fmt.Println("New sell Order:", string(bytes))

				database.FindCycleByIdAndUpdate(idString, "status", "sell")
			}
		} else if status == "sell" {
			order, _ := client.GetOrderById((buyId).(string))
			isFilled := client.IsFilled(string(order))

			if !isFilled {
				fmt.Println("Order Not filled")
			} else {
				database.FindCycleByIdAndUpdate(idString, "status", "completed")
				color.Green("Cycle successfully completed")
			}
		} else {
			log.Fatal("Unsupported status:", status)
		}
	}
}
