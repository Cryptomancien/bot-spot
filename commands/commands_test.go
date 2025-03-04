package commands

import (
	"fmt"
	"testing"
)

func TestGetClientByExchange(t *testing.T) {
	client := GetClientByExchange()
	fmt.Println(client)
}
