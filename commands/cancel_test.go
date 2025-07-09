package commands

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func TestCancel(t *testing.T) {
	err := godotenv.Load("../bot.conf")
	if err != nil {
		t.Fatalf("Error loading ../bot.conf: %v", err)
	}

	orderID := os.Getenv("ORDER_ID")
	if orderID == "" {
		t.Fatal("Missing environment variable: ORDER_ID")
	}

	client := GetClientByExchange()
	client.CheckConnection()

	cancelOrder, err := client.CancelOrder(orderID)
	if err != nil {
		log.Fatal("Cannot cancel order: ", string(cancelOrder))
	}
}
