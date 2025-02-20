package mexc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

func (c *Client) SetBaseURL(url string) {
	c.BaseURL = url
}

func (c *Client) CheckConnection() {
	resp, err := http.Get(c.BaseURL + "/api/v3/ping")
	if err != nil {
		log.Fatalf("Failed to connect to MEXC: %v", err)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("MEXC API returned non-200 status: %d", resp.StatusCode)
	}

	fmt.Println("Connected to MEXC API successfully")
}

func (c *Client) GetBalanceUSDT() float64 {
	color.Blue("Checking USDT balance...")
	endpoint := "/api/v3/account"
	timestamp := time.Now().UnixMilli()

	queryString := fmt.Sprintf("timestamp=%d", timestamp)
	h := hmac.New(sha256.New, []byte(c.APISecret))
	h.Write([]byte(queryString))
	signature := hex.EncodeToString(h.Sum(nil))

	fullURL := fmt.Sprintf("%s%s?%s&signature=%s", c.BaseURL, endpoint, queryString, signature)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Println("Error creating request:", err)
	}

	req.Header.Set("X-MEXC-APIKEY", c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making request:", err)
		os.Exit(0)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		os.Exit(0)
	}

	balances, _, _, err := jsonparser.Get(body, "balances")
	if err != nil {
		log.Println("Error getting balances:", err)
		os.Exit(0)
	}

	var freeFloat float64
	_, err = jsonparser.ArrayEach(balances, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		asset, _ := jsonparser.GetString(value, "asset")
		if asset == "USDT" {
			freeStr, _ := jsonparser.GetString(value, "free")
			free, _ := strconv.ParseFloat(freeStr, 64)
			if err != nil {
				log.Println("Error converting free balance:", err)
				os.Exit(0)
			}

			freeFloat = free
		}
	})

	return freeFloat
}

func (c *Client) GetLastPriceBTC() float64 {
	fmt.Println("Getting last price...")

	endpoint := "/api/v3/ticker/price"
	symbol := "BTCUSDT"

	fullURL := fmt.Sprintf("%s%s?symbol=%s", c.BaseURL, endpoint, symbol)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	req.Header.Set("X-MEXC-APIKEY", c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
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

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Fatal("Error converting price:", err)
		return 0
	}

	return price
}

func (c *Client) CreateOrder(side string, price, quantity float64) string {
	endpoint := "/api/v3/order"
	timestamp := time.Now().UnixMilli()

	queryString := fmt.Sprintf("symbol=BTCUSDT&side=%s&type=LIMIT&timeInForce=GTC&quantity=%f&price=%f&timestamp=%d",
		side, quantity, price, timestamp)

	h := hmac.New(sha256.New, []byte(c.APISecret))
	h.Write([]byte(queryString))
	signature := hex.EncodeToString(h.Sum(nil))

	fullURL := fmt.Sprintf("%s%s?%s&signature=%s", c.BaseURL, endpoint, queryString, signature)

	req, err := http.NewRequest("POST", fullURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("X-MEXC-APIKEY", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println("Error: Received non-OK HTTP status code:", resp.StatusCode)
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Order Response:", string(body))

	var orderResponse struct {
		OrderID int `json:"orderId"`
	}
	if err := json.Unmarshal(body, &orderResponse); err != nil {
		log.Fatal(err)
	}

	return strconv.Itoa(orderResponse.OrderID)
}
