package commands

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"main/database"
	"os"
	"strconv"
	"strings"
)

func Cancel() {
	lastArg := GetLastArg()
	delimiter := "="

	_, numberString, found := strings.Cut(lastArg, delimiter)
	if !found {
		log.Fatal("Could not find delimiter")
	}

	idInt, err := strconv.Atoi(numberString)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Cancelling", idInt)

	document := database.GetByIdInt(idInt)

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
	} else {
		color.Red("Unknown status")
		os.Exit(0)
	}

	client := GetClientByExchange()

	res, err := client.CancelOrder(orderIdToCancel)
	if err != nil {
		return
	}
	fmt.Println(string(res))

	database.DeleteByIdInt(int32(idInt))

	color.Green("Cycle %d successfully canceled", idInt)
}
