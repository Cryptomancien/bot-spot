package mexc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
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

// Generates HMAC SHA256 signature for a signed request
func (c *Client) signRequest(queryString string) string {
	h := hmac.New(sha256.New, []byte(c.APISecret))
	h.Write([]byte(queryString))
	return hex.EncodeToString(h.Sum(nil))
}

// Sends an HTTP request and returns the response body
func (c *Client) sendRequest(method, endpoint, queryString string) ([]byte, error) {
	fullURL := fmt.Sprintf("%s%s?%s", c.BaseURL, endpoint, queryString)

	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-MEXC-APIKEY", c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("error: received non-OK HTTP status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) CheckConnection() {
	_, err := c.sendRequest("GET", "/api/v3/ping", "")
	if err != nil {
		log.Fatalf("Failed to connect to MEXC: %v", err)
	}

	fmt.Println("Connected to MEXC API successfully")
}

func (c *Client) GetBalanceUSDT() float64 {
	fmt.Println("Checking USDT balance...")

	timestamp := time.Now().UnixMilli()
	queryString := fmt.Sprintf("timestamp=%d", timestamp)
	signature := c.signRequest(queryString)
	signedQuery := fmt.Sprintf("%s&signature=%s", queryString, signature)

	body, err := c.sendRequest("GET", "/api/v3/account", signedQuery)
	if err != nil {
		log.Fatalf("Error fetching balance: %v", err)
	}

	balances, _, _, err := jsonparser.Get(body, "balances")
	if err != nil {
		log.Fatalf("Error getting balances: %v", err)
	}

	var freeFloat float64
	_, err = jsonparser.ArrayEach(balances, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		asset, _ := jsonparser.GetString(value, "asset")
		if asset == "USDT" {
			freeStr, _ := jsonparser.GetString(value, "free")
			free, _ := strconv.ParseFloat(freeStr, 64)
			freeFloat = free
		}
	})

	return freeFloat
}

func (c *Client) GetLastPriceBTC() float64 {
	fmt.Println("Fetching last BTC price...")

	queryString := "symbol=BTCUSDT"
	body, err := c.sendRequest("GET", "/api/v3/ticker/price", queryString)
	if err != nil {
		log.Fatalf("Error fetching BTC price: %v", err)
	}

	priceStr, err := jsonparser.GetString(body, "price")
	if err != nil {
		log.Fatalf("Error extracting price: %v", err)
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Fatalf("Error converting price: %v", err)
	}

	return price
}

func (c *Client) CreateOrder(side string, price, quantity float64) string {
	timestamp := time.Now().UnixMilli()

	queryString := fmt.Sprintf(
		"symbol=BTCUSDT&side=%s&type=LIMIT&timeInForce=GTC&quantity=%.6f&price=%.2f&timestamp=%d",
		side, quantity, price, timestamp,
	)

	signature := c.signRequest(queryString)
	signedQuery := fmt.Sprintf("%s&signature=%s", queryString, signature)

	body, err := c.sendRequest("POST", "/api/v3/order", signedQuery)
	if err != nil {
		log.Fatalf("Error placing order: %v", err)
	}

	var orderResponse struct {
		OrderID int `json:"orderId"`
	}
	if err := json.Unmarshal(body, &orderResponse); err != nil {
		log.Fatalf("Error parsing order response: %v", err)
	}

	return strconv.Itoa(orderResponse.OrderID)
}
