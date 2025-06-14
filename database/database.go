package database

import (
	"errors"
	"github.com/ostafen/clover"
	"log"
	"os"
	"path/filepath"
)

var db *clover.DB

const CollectionName = "cycles"

func InitDatabase() {
	databasePath := GetDatabasePath()
	var err error
	db, err = clover.Open(databasePath)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	collectionAlreadyExists, _ := db.HasCollection(CollectionName)
	if !collectionAlreadyExists {
		if err := db.CreateCollection(CollectionName); err != nil {
			log.Fatal(err)
		}
	}
}

func GetDatabasePath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	rootDir := filepath.Join(dirname, "cryptomancien")
	databasePath := filepath.Join(rootDir, "bot-v3")

	if _, err := os.Stat(rootDir); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(rootDir, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat(databasePath); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(databasePath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	return databasePath
}

type Status string

type Cycle struct {
	IdInt     int32
	Exchange  string
	Status    string
	Quantity  float64
	BuyPrice  float64
	BuyId     string
	SellPrice float64
	SellId    string
}

func List() []*clover.Document {
	docs, err := db.Query(CollectionName).Sort(clover.SortOption{Field: "idInt", Direction: -1}).FindAll()
	if err != nil {
		log.Fatal(err)
	}
	return docs
}

func PrepareIdInt() int32 {
	count, err := db.Query(CollectionName).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		return 1
	}

	lastDoc, err := db.Query(CollectionName).Sort(clover.SortOption{
		Field:     "idInt",
		Direction: -1,
	}).Limit(1).FindFirst()
	if err != nil {
		log.Fatal(err)
	}

	lastId := (lastDoc.Get("idInt")).(int64)
	return int32(lastId + 1)
}

func NewCycle(cycle *Cycle) string {
	idInt := PrepareIdInt()

	doc := clover.NewDocument()
	doc.Set("idInt", idInt)
	doc.Set("exchange", cycle.Exchange)
	doc.Set("status", cycle.Status)
	doc.Set("quantity", cycle.Quantity)
	doc.Set("buyPrice", cycle.BuyPrice)
	doc.Set("buyId", cycle.BuyId)
	doc.Set("sellPrice", cycle.SellPrice)
	doc.Set("sellId", cycle.SellId)

	docId, err := db.InsertOne(CollectionName, doc)
	if err != nil {
		log.Fatal(err)
	}

	return docId
}

func GetById(id string) *clover.Document {
	document, err := db.Query(CollectionName).FindById(id)
	if err != nil {
		log.Fatal(err)
	}
	return document
}

func GetByIdInt(id int) *clover.Document {
	document, err := db.Query(CollectionName).Where(clover.Field("idInt").Eq(id)).FindFirst()
	if err != nil {
		log.Println(err)
	}
	return document
}

func DeleteById(id string) {
	err := db.Query(CollectionName).DeleteById(id)
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteByIdInt(idInt int32) {
	err := db.Query(CollectionName).Where(clover.Field("idInt").Eq(idInt)).Delete()
	if err != nil {
		log.Fatal(err)
	}
}

func ListPerPage(page, perPage int) []*clover.Document {
	skip := (page - 1) * perPage
	docs, err := db.Query(CollectionName).
		Sort(clover.SortOption{Field: "idInt", Direction: -1}).
		Skip(skip).
		Limit(perPage).
		FindAll()
	if err != nil {
		log.Fatal(err)
	}
	return docs
}

func FindCycleByIdAndUpdate(id, field string, value interface{}) {
	err := db.Query(CollectionName).UpdateById(id, map[string]interface{}{field: value})
	if err != nil {
		log.Fatal(err)
	}
}
