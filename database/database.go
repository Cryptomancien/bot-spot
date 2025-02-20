package database

import (
	"errors"
	"fmt"
	"github.com/ostafen/clover"
	"log"
	"os"
	"path/filepath"
)

const CollectionName = "cycles"

func InitDatabase() {
	databasePath := GetDatabasePath()
	db, _ := clover.Open(databasePath)
	collectionAlreadyExists, _ := db.HasCollection(CollectionName)
	if !collectionAlreadyExists {
		err := db.CreateCollection(CollectionName)
		if err != nil {
			log.Fatal(err)
		}
	}
	err := db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func GetDatabasePath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	databasePath := filepath.Join(dirname, "cryptomancien/bot-v3")
	if _, err := os.Stat(databasePath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(databasePath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	return databasePath
}

func GetDB() *clover.DB {
	databasePath := GetDatabasePath()
	db, err := clover.Open(databasePath)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	return db
}

type Cycle struct {
	Exchange  string
	Status    string
	Quantity  float32
	BuyPrice  float32
	BuyId     string
	SellPrice float32
	SellId    string
}

func List() []*clover.Document {
	db := GetDB()
	defer func(db *clover.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	docs, err := db.Query(CollectionName).FindAll()
	if err != nil {
		log.Fatal(err)
	}

	return docs
}

func NewCycle(cycle *Cycle) {
	exchange := cycle.Exchange
	status := cycle.Status
	quantity := cycle.Quantity
	buyPrice := cycle.BuyPrice
	buyId := cycle.BuyId
	sellPrice := cycle.SellPrice
	sellId := cycle.SellId

	doc := clover.NewDocument()
	doc.Set("exchange", exchange)
	doc.Set("status", status)
	doc.Set("quantity", quantity)
	doc.Set("buyPrice", buyPrice)
	doc.Set("buyId", buyId)
	doc.Set("sellPrice", sellPrice)
	doc.Set("sellId", sellId)

	db := GetDB()
	docId, err := db.InsertOne(CollectionName, doc)
	if err != nil {
		return
	}

	fmt.Println("docId:", docId)

	defer func(db *clover.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)
}

func GetById(id string) *clover.Document {
	db := GetDB()

	document, err := db.Query(CollectionName).FindById(id)
	if err != nil {
		log.Fatal(err)
	}

	return document
}

func DeleteById(id string) {
	db := GetDB()
	err := db.Query(CollectionName).DeleteById(id)
	if err != nil {
		log.Fatal(err)
	}
}
