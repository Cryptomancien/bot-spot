package database

import (
	"errors"
	"github.com/ostafen/clover"
	"log"
	"os"
	"path/filepath"
)

const CollectionName = "cycles"

func InitDatabase() (*clover.DB, error) {
	databasePath := GetDatabasePath()
	db, _ := clover.Open(databasePath)
	collectionAlreadyExists, _ := db.HasCollection(CollectionName)
	if !collectionAlreadyExists {
		err := db.CreateCollection(CollectionName)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer func(db *clover.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Error db.Close: %s", err)
		}
	}(db)

	return db, nil
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

func GetDB() *clover.DB {
	databasePath := GetDatabasePath()
	db, err := clover.Open(databasePath)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	return db
}

type Status string

type Cycle struct {
	IdInt     int32
	Exchange  string
	Status    string // buy sell completed
	Quantity  float64
	BuyPrice  float64
	BuyId     string
	SellPrice float64
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

	docs, err := db.Query(CollectionName).Sort(clover.SortOption{Field: "idInt", Direction: -1}).FindAll()
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatal(err)
	}

	return docs
}

func PrepareIdInt() int32 {
	db := GetDB()

	defer func(db *clover.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	count, err := db.Query(CollectionName).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		return 1
	}

	lastDoc, _ := db.Query(CollectionName).Sort(clover.SortOption{
		Field:     "idInt",
		Direction: -1,
	}).Limit(1).FindFirst()

	lastId := (lastDoc.Get("idInt")).(int64)
	nextId := lastId + 1

	return int32(nextId)
}

func NewCycle(cycle *Cycle) string {
	idInt := PrepareIdInt()
	exchange := cycle.Exchange
	status := cycle.Status
	quantity := cycle.Quantity
	buyPrice := cycle.BuyPrice
	buyId := cycle.BuyId
	sellPrice := cycle.SellPrice
	sellId := cycle.SellId

	doc := clover.NewDocument()
	doc.Set("idInt", idInt)
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
		log.Fatal(err)
	}

	defer func(db *clover.DB) {
		err := db.Close()
		if err != nil {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Fatal(err)
		}
	}(db)

	return docId
}

func GetById(id string) *clover.Document {
	db := GetDB()

	defer func(db *clover.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	document, err := db.Query(CollectionName).FindById(id)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatal(err)
	}
	return document
}

func GetByIdInt(id int) *clover.Document {
	db := GetDB()

	defer func(db *clover.DB) {
		err := db.Close()
		if err != nil {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Fatal(err)
		}
	}(db)

	document, err := db.Query(CollectionName).Where(clover.Field("idInt").Eq(id)).FindFirst()
	if err != nil {
		log.Println(err)
	}

	return document
}

func DeleteById(id string) {
	db := GetDB()

	defer func(db *clover.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	err := db.Query(CollectionName).DeleteById(id)
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteByIdInt(idInt int32) {
	db := GetDB()

	defer func(db *clover.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	err := db.
		Query(CollectionName).
		Where(clover.Field("idInt").
			Eq(idInt)).
		Delete()
	if err != nil {
		log.Fatal(err)
	}
}

func ListPerPage(page, perPage int) []*clover.Document {
	db := GetDB()
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	skip := (page - 1) * perPage
	docs, _ := db.Query(CollectionName).Sort(clover.SortOption{
		Field:     "idInt",
		Direction: -1,
	}).
		Skip(skip).
		Limit(perPage).
		FindAll()

	return docs
}

func FindCycleByIdAndUpdate(id, field string, value interface{}) {
	db := GetDB()
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	err := db.Query(CollectionName).UpdateById(id, map[string]interface{}{field: value})
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatal(err)
	}
}
