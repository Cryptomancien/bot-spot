package exchange

import (
	"context"
	"fmt"
	"github.com/bogdankorobka/mexc-golang-sdk"
	mexchttp "github.com/bogdankorobka/mexc-golang-sdk/http"
	"github.com/fatih/color"
	"log"
	"net/http"

	"os"
)

func CheckConnection() {
	fmt.Println("Check connection...")
	exchange := os.Getenv("EXCHANGE")
	if exchange == "" {
		color.Red("You need to set EXCHANGE in your environment")
	}

	switch exchange {
	case "MEXC":
		apiKey := os.Getenv("MEXC_API_KEY")
		apiSecret := os.Getenv("MEXC_API_SECRET")

		client := mexchttp.NewClient(apiKey, apiSecret, &http.Client{})
		if client == nil {
			log.Fatal("Error creating MEXC client")
		}
		ctx := context.Background()
		rest, err := mexc.NewRest(ctx, client)
		if err != nil {
			log.Fatal("Error creating MEXC REST client:", err)
		}

		ping, _ := rest.MarketService.Ping(ctx)
		if ping == "{}" {
			color.Green("Connection OK")
		} else {
			os.Exit(0)
		}
	}
}

func CheckBalanceUSDT() {
	fmt.Println("Check balance USDT...")
	exchange := os.Getenv("EXCHANGE")

	switch exchange {
	case "MEXC":
		apiKey := os.Getenv("MEXC_API_KEY")
		apiSecret := os.Getenv("MEXC_API_SECRET")

		client := mexchttp.NewClient(apiKey, apiSecret, &http.Client{})
		if client == nil {
			log.Fatal("Error creating MEXC client")
		}
		ctx := context.Background()
		_, err := mexc.NewRest(ctx, client)
		if err != nil {
			log.Fatal("Error creating MEXC REST client:", err)
		}

		account, err := client.SendRequest(ctx, "GET", "/api/v3/accounts", nil)
		if err != nil {

		}
		if len(account) == 0 {
			color.Red("No account found")
			os.Exit(0)
		}
	}
}
