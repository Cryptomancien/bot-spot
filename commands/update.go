package commands

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/fatih/color"
	"log"
	"main/database"
	"main/tools"
	"os"
	"strconv"
)

func Update() {
	MainMiddleware()
	client := GetClientByExchange()

	lastPrice := client.GetLastPriceBTC()
	//fmt.Println("Last price:", lastPrice)

	docs := database.List()
	for _, doc := range docs {

		id := doc.Get("_id")
		idInt := doc.Get("idInt")
		idString := id.(string)

		status := doc.Get("status")
		quantity := doc.Get("quantity")
		buyPrice := doc.Get("buyPrice")
		buyId := doc.Get("buyId")
		sellPrice := doc.Get("sellPrice")
		sellId := doc.Get("sellId")

		if status == "buy" {
			order, _ := client.GetOrderById((buyId).(string))
			isFilled, err := client.IsFilled(string(order))
			if err != nil {
				color.Red("Error found on cycle %d (don't worry)\n", idInt)
				color.Yellow("Try to remove it")

				fmt.Printf("go run . -c %d\n", idInt)
				fmt.Printf("go run . -cl %d %d\n", idInt, idInt)
				os.Exit(0)
			}

			if !isFilled {
				fmt.Printf("%s %s %s\n",
					color.YellowString("%d", idInt),
					color.CyanString("Order Buy  still active -"),
					color.WhiteString("%s", buyId),
				)
			} else {
				fmt.Printf("%s %s\n",
					color.YellowString("%d", idInt),
					color.GreenString("Order Buy filled"),
				)

				// Check sell price > last price
				fmt.Println("SELL PRICE", sellPrice)
				fmt.Println("id ", idInt)
				sellPrice := (sellPrice).(float64)

				if lastPrice > sellPrice {
					const offset = 200.0
					newSellPrice := sellPrice + offset
					fmt.Println("New sell price: ", newSellPrice)

					database.FindCycleByIdAndUpdate(idString, "sellPrice", newSellPrice)
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
					log.Fatalf("Error creating order: %v", err)
					return
				}
				orderId, _, _, err := jsonparser.Get(bytes, "orderId")

				if err != nil {
					log.Printf("Cycle id %s", idInt)
					log.Printf("Failed to parse orderId: %v", err)
					log.Fatal(err)
				}

				fmt.Printf("%s %s %s\n",
					color.YellowString("%d", idInt),
					color.CyanString("New sell Order -"),
					color.WhiteString("%s", string(bytes)),
				)

				database.FindCycleByIdAndUpdate(idString, "status", "sell")
				database.FindCycleByIdAndUpdate(idString, "sellId", string(orderId))
			}
		} else if status == "sell" {
			order, _ := client.GetOrderById((sellId).(string))
			isFilled, err := client.IsFilled(string(order))
			if err != nil {
				color.Red("Error found on cycle %d (don't worry)\n", idInt)
				color.Yellow("Try to remove it")

				fmt.Printf("go run . -c %d\n", idInt)
				fmt.Printf("go run . -cl %d %d\n", idInt, idInt)
				os.Exit(0)
			}

			if !isFilled {
				fmt.Printf("%s %s %s\n",
					color.YellowString("%d", idInt),
					color.CyanString("Order Sell still active -"),
					color.WhiteString("%s", sellId),
				)
			} else {
				database.FindCycleByIdAndUpdate(idString, "status", "completed")

				// Calc Percent
				totalBuyUSD := buyPrice.(float64) * quantity.(float64)
				totalSellUSD := sellPrice.(float64) * quantity.(float64)
				percent := (totalSellUSD - totalBuyUSD) / totalBuyUSD * 100

				fmt.Printf("%s %s (Gain: %.2f%%)\n",
					color.YellowString("%d", idInt),
					color.GreenString("Cycle successfully completed"),
					percent,
				)

				if os.Getenv("TELEGRAM") == "1" {
					var message = ""
					message += fmt.Sprintf("✅ Cycle %d completed \n", idInt)
					message += fmt.Sprintf("📉 Buy Price: %.2f \n", buyPrice)
					message += fmt.Sprintf("📈 Sell Price: %.2f \n", sellPrice)
					message += fmt.Sprintf("💰 Gain: $ %.2f \n", totalSellUSD-totalBuyUSD)
					tools.Telegram(message)
				}
			}
		}
	}
	Log("Update complete")
}
