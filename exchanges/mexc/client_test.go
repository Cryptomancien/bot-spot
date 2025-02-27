package mexc

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"testing"
)

var client *Client

func TestMain(m *testing.M) {
	// TODO make config folder
	_ = godotenv.Load("../../bot.conf")

	client = NewClient()
	client.SetBaseURL("https://api.mexc.co")

	os.Exit(m.Run())
}

func TestCheckConnection(t *testing.T) {
	client.CheckConnection()
}

func TestClient_GetBalanceUSD(t *testing.T) {
	balance := client.GetBalanceUSD()
	fmt.Println("balance:", balance)
}
