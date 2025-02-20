package database

import (
	"fmt"
	"testing"
)

func TestNewCycle(t *testing.T) {
	cycle := Cycle{
		Exchange:  "mexc",
		Status:    "buy",
		Quantity:  0.004,
		BuyPrice:  98000,
		BuyId:     "231123",
		SellPrice: 99000,
		SellId:    "",
	}

	NewCycle(&cycle)
}

func TestList(t *testing.T) {
	documents := List()
	for _, document := range documents {
		fmt.Println(document)
	}
}

func TestGetById(t *testing.T) {
	id := "337bded7-7c85-46ff-a5d1-9c6314dd9826"
	doc := GetById(id)
	fmt.Println(doc)
}

func TestDeleteById(t *testing.T) {
	id := "6a3743ae-bee8-46da-9905-81b8ae261a72"
	DeleteById(id)
}
