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
	id := "ef26b72f-e034-45db-95e2-543ad40f229f"
	DeleteById(id)
}

func TestPrepareIdInt(t *testing.T) {
	id := PrepareIdInt()
	fmt.Println(id)
}
