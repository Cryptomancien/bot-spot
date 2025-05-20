package kucoin

import (
	"os"
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
		BaseURL:   "https://api.kucoin.com/api",
	}
}

func (c *Client) sendRequest(method, endpoint, queryString string) /*([]byte, error)*/ {

}

func (c *Client) SetBaseURL(url string) {
	c.BaseURL = url
}

func (c *Client) CheckConnection() {

}

func (c *Client) GetBalanceUSD() /*float64*/ {

}

func (c *Client) GetLastPriceBTC() /*float64*/ {

}

func CreateOrder() {

}

func GetOrderById() {

}

func IsFilled() {

}

func CancelOrder() {

}
