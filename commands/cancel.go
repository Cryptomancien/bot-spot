package commands

import (
	// "fmt"
	// "github.com/fatih/color"
	// "log"
	// "main/database"
	// "os"
	// "strconv"
	// "strings"
	"fmt"
	"github.com/fatih/color"
	"log"
	"main/database"
	"os"
	"strconv"
	"strings"
)

func Cancel() {
	if len(os.Args) < 3 {
		color.Red("Id required")
		color.Cyan("go run . -c 34")
		os.Exit(1)
	}

	lastArg := os.Args[2]

	idInt, err := strconv.Atoi(lastArg)
	if err != nil {
		log.Fatal(err)
	}

	color.Yellow("Cancelling %d", idInt)

	document := database.GetByIdInt(idInt)
	if document == nil {
		color.Red("No document found")
		return
	}

	status := document.Get("status").(string)
	exchange := document.Get("exchange").(string)
	exchange = strings.ToUpper(exchange)

	if status == "completed" {
		color.Red("Can't cancel completed cycle, only 'buy' or 'sell' is supported")
		os.Exit(0)
	}

	var orderIdToCancel string
	if status == "buy" {
		orderIdToCancel = (document.Get("buyId")).(string)
	} else if status == "sell" {
		orderIdToCancel = (document.Get("sellId")).(string)

		if orderIdToCancel == "" {
			color.Cyan("No order found for cancelling...deleting")
			database.DeleteByIdInt(int32(idInt))
			color.Green("Cycle %d successfully canceled", idInt)
			os.Exit(0)
		}
	} else {
		color.Red("Unknown status")
		os.Exit(0)
	}

	client := GetClientByExchange()

	res, err := client.CancelOrder(orderIdToCancel)
	if err != nil {
		fmt.Println(string(res))
		return
	}

	database.DeleteByIdInt(int32(idInt))

	color.Green("Cycle %d successfully canceled", idInt)
}
