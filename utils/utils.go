package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
)

const ConfigFilename = "bot.conf"

func CreateConfigFileIfNotExists() {
	if _, err := os.Stat(ConfigFilename); errors.Is(err, os.ErrNotExist) {
		content := []byte("CUSTOMER_ID=\n\nEXCHANGE=\n\nMEXC_PUBLIC=\nMEXC_PRIVATE=\n\nBUY_OFFSET=-1000\nSELL_OFFSET=1000\n\nPERCENT=6")
		err := os.WriteFile(ConfigFilename, content, 0644)
		if err != nil {
			panic(err)
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
