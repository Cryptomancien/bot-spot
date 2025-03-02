package database

import (
	"fmt"
	"github.com/ostafen/clover"
	"log"
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

func TestDeleteByIdInt(t *testing.T) {
	//for i := range 13 {
	//}
	id := 42
	DeleteByIdInt(int32(id))
}

func TestPrepareIdInt(t *testing.T) {
	id := PrepareIdInt()
	fmt.Println(id)
}

func TestFindCycleByIdAndUpdate(t *testing.T) {
	id := "191edf92-0d31-48cc-91d4-af0a15b9428d"
	db := GetDB()

	defer func(db *clover.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	err := db.Query(CollectionName).UpdateById(id, map[string]interface{}{"sellPrice": 86156.31})
	if err != nil {
		log.Fatal(err)
	}
}
