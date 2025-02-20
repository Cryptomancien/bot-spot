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
		BaseURL:   "https://api.mexc.com",
	}
}

func (c *Client) SetBaseURL(url string) {
	c.BaseURL = url
}

func (c *Client) sign(query string) string {
	h := hmac.New(sha256.New, []byte(c.APISecret))
	h.Write([]byte(query))
	return hex.EncodeToString(h.Sum(nil))
}

func (c *Client) request(method, endpoint, queryString string) ([]byte, error) {
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) CheckConnection() {
	color.Blue("Checking connection...")
	if _, err := c.request("GET", "/api/v3/ping", ""); err != nil {
		log.Fatal("Connection failed:", err)
	}
	color.Green("Connection OK")
}

func (c *Client) GetBalanceUSDT() float32 {
	color.Blue("Checking USDT balance...")
	timestamp := time.Now().UnixMilli()
	signature := c.sign(fmt.Sprintf("timestamp=%d", timestamp))
	body, err := c.request("GET", "/api/v3/account", fmt.Sprintf("timestamp=%d&signature=%s", timestamp, signature))
	if err != nil {
		log.Fatal("Request failed:", err)
	}

	balances, _, _, err := jsonparser.Get(body, "balances")
	if err != nil {
		log.Fatal("Error parsing balances:", err)
	}

	var balance float32
	_, err = jsonparser.ArrayEach(balances, func(value []byte, _ jsonparser.ValueType, _ int, _ error) {
		if asset, _ := jsonparser.GetString(value, "asset"); asset == "USDT" {
			freeStr, _ := jsonparser.GetString(value, "free")
			free, _ := strconv.ParseFloat(freeStr, 32)
			balance = float32(free)
		}
	})
	if err != nil {
		return 0
	}

	return balance
}

func (c *Client) GetLastPriceBTC() float32 {
	fmt.Println("Getting last price...")
	body, err := c.request("GET", "/api/v3/ticker/price", "symbol=BTCUSDT")
	if err != nil {
		log.Fatal("Request failed:", err)
	}

	priceStr, err := jsonparser.GetString(body, "price")
	if err != nil {
		log.Fatal("Error extracting price:", err)
	}

	price, _ := strconv.ParseFloat(priceStr, 32)
	return float32(price)
}

func (c *Client) CreateOrder(side string, price, quantity float32) string {
	fmt.Println("Creating order...")
	timestamp := time.Now().UnixMilli()
	query := fmt.Sprintf("symbol=BTCUSDT&side=%s&type=LIMIT&timeInForce=GTC&quantity=%f&price=%f&timestamp=%d",
		side, quantity, price, timestamp)
	signature := c.sign(query)

	body, err := c.request("POST", "/api/v3/order", fmt.Sprintf("%s&signature=%s", query, signature))
	if err != nil {
		log.Fatal("Request failed:", err)
	}

	var orderResponse struct {
		OrderID int `json:"orderId"`
	}
	if err := json.Unmarshal(body, &orderResponse); err != nil {
		log.Fatal("Error parsing response:", err)
	}

	return strconv.Itoa(orderResponse.OrderID)
}
