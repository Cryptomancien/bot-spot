package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"io"
	"log"
	"main/exchanges/mexc"
	"net/http"
	"os"
	"strings"
)

type ExchangeClient interface {
	CheckConnection()
	GetBalanceUSD() float64
	GetLastPriceBTC() float64
	SetBaseURL(url string)
	CreateOrder(side, price, quantity string) ([]byte, error)
	GetOrderById(id string) ([]byte, error)
	IsFilled(id string) bool
	CancelOrder(orderID string) ([]byte, error)
}

const ConfigFilename = "bot.conf"

func CreateConfigFileIfNotExists() {
	if _, err := os.Stat(ConfigFilename); errors.Is(err, os.ErrNotExist) {
		pathConfTemplate := fmt.Sprintf("commands/misc/%s.example", ConfigFilename)

		content, err := os.ReadFile(pathConfTemplate)
		if err != nil {
			content = []byte("CUSTOMER_ID=\n\nEXCHANGE=\n\nMEXC_PUBLIC=\nMEXC_PRIVATE=\n\nBUY_OFFSET=-1000\nSELL_OFFSET=1000\n\nPERCENT=6\n\nAUTO_INTERVAL_NEW=60\nAUTO_INTERVAL_UPDATE=10")
			err := os.WriteFile(ConfigFilename, content, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = os.WriteFile(ConfigFilename, content, 0644)
		if err != nil {
			log.Fatal(err)
		}
		color.Green("Config file created: " + ConfigFilename)
	}

}

func LoadDotEnv() {
	err := godotenv.Load(ConfigFilename)
	if err != nil {
		log.Fatal("Error loading bot.conf")
	}
}

func GetLastArg() string {
	args := os.Args
	argsLen := len(args)
	lastArg := args[argsLen-1]
	return lastArg
}

func CheckPremium() {
	color.Blue("Checking subscription before...")
	var customerId string = os.Getenv("CUSTOMER_ID")

	if customerId == "" {
		color.Red("You need to set CUSTOMER_ID in bot.conf")
		os.Exit(0)
	}

	url := "https://validator.cryptomancien.com"
	body, err := json.Marshal(
		map[string]string{
			"CUSTOMER_ID": os.Getenv("CUSTOMER_ID"),
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	statusCode := resp.StatusCode
	if statusCode != 200 {
		color.Red("CUSTOMER_ID not in bot.conf or subscription expired")
		color.Red("Go to https://cryptomancien.com -> Space -> Trading bots and get your CUSTOMER_ID")
		color.Red("Then fill it in your config file bot.conf")
		os.Exit(0)
	}

	color.Green("Subscription OK")
	fmt.Println("")
}

func GetClientByExchange(exchangeArg ...string) ExchangeClient {

	var ex string
	if len(exchangeArg) > 0 {
		ex = exchangeArg[0]
	} else {
		ex = os.Getenv("EXCHANGE")
	}
	ex = strings.ToUpper(ex)

	var client ExchangeClient

	switch ex {
	case "MEXC":
		client = mexc.NewClient()
		client.SetBaseURL("https://api.mexc.co")
	default:
		fmt.Println("Unsupported exchange:", ex)
		os.Exit(0)
	}

	return client
}
