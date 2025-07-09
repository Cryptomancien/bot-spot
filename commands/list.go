package commands

import (
	"encoding/json"
	"log"
)

func List() {
	client := GetClientByExchange()

	client.CheckConnection()

	orders, err := client.GetOpenOrders()
	if err != nil {
		panic(err)
	}

	var prettyJSON interface{}
	err = json.Unmarshal(orders, &prettyJSON)
	if err != nil {
		log.Fatalf("Error unmarshalling orders JSON: %v", err)
	}

	indentedOrders, err := json.MarshalIndent(prettyJSON, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling orders to indented JSON: %v", err)
	}

	log.Println("Open Orders:\n", string(indentedOrders))
}
