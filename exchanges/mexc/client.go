package mexc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/fatih/color"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Client struct {
	APIKey    string
	APISecret string
	BaseURL   string
}

func NewClient() *Client {
	return &Client{
		APIKey:    os.Getenv("MEXC_API_KEY"),
		APISecret: os.Getenv("MEXC_SECRET_KEY"),
		BaseURL:   "https://api.mexc.co",
	}
}

func (c *Client) CheckConnection() {
	color.Blue("Checking connection...")

	endpoint := "/api/v3/ping"
	fullURL := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(0)
	}

	// Perform the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		os.Exit(0)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing body:", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		os.Exit(0)
	}

	color.Green("Connection OK")
	fmt.Println("")
}

func (c *Client) GetBalanceUSDT() float32 {
	color.Blue("Checking USDT balance...")
	endpoint := "/api/v3/account"
	timestamp := time.Now().UnixMilli()

	// Construct the query string
	queryString := fmt.Sprintf("timestamp=%d", timestamp)

	// Generate HMAC-SHA256 signature
	h := hmac.New(sha256.New, []byte(c.APISecret))
	h.Write([]byte(queryString))
	signature := hex.EncodeToString(h.Sum(nil))

	// Build the full request URL
	fullURL := fmt.Sprintf("%s%s?%s&signature=%s", c.BaseURL, endpoint, queryString, signature)

	// Create the HTTP request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	// Set request headers
	req.Header.Set("X-MEXC-APIKEY", c.APIKey)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making request:", err)
		os.Exit(0)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing body:", err)
		}
	}(resp.Body)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		os.Exit(0)
	}
	// {"makerCommission":null,"takerCommission":null,"buyerCommission":null,"sellerCommission":null,"canTrade":true,"canWithdraw":true,"canDeposit":true,"updateTime":null,"accountType":"SPOT","balances":[{"asset":"USDT","free":"25789.44","locked":"0"},{"asset":"SOL","free":"0.009678398","locked":"0"}],"permissions":["SPOT"]}

	balances, _, _, err := jsonparser.Get(body, "balances")
	if err != nil {
		log.Println("Error getting balances:", err)
		os.Exit(0)
	}

	var freeFloat float32
	_, err = jsonparser.ArrayEach(balances, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		asset, _ := jsonparser.GetString(value, "asset")
		if asset == "USDT" {
			freeStr, _ := jsonparser.GetString(value, "free")
			free, _ := strconv.ParseFloat(freeStr, 32)
			if err != nil {
				fmt.Println("Error converting free balance:", err)
				os.Exit(0)
			}

			freeFloat = float32(free)
		}
	})

	return freeFloat
}

func (c *Client) GetLastPriceBTC() float32 {
	fmt.Println("Getting last price...")

	endpoint := "/api/v3/ticker/price"
	symbol := "BTCUSDT"

	// Build full URL
	fullURL := fmt.Sprintf("%s%s?symbol=%s", c.BaseURL, endpoint, symbol)

	// Create HTTP request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	// Set headers if required
	req.Header.Set("X-MEXC-APIKEY", c.APIKey)

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing body:", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	priceStr, err := jsonparser.GetString(body, "price")
	if err != nil {
		log.Fatal("Error extracting price:", err)
		return 0
	}

	price, err := strconv.ParseFloat(priceStr, 32)
	if err != nil {
		log.Fatal("Error converting price:", err)
		return 0
	}

	return float32(price)
}
