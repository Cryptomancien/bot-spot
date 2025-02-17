package debug

import (
	"fmt"
	"github.com/ostafen/clover"
)

func TestDebug() {

	db, _ := clover.Open("cycles.db")

	//err := db.CreateCollection("cycles")
	//if err != nil {
	//	return
	//}
	//

	cycle := &struct {
		Id        string `clover:"_id"`
		Exchange  string `clover:"exchange"`
		Status    string `clover:"status"`
		Quantity  string `clover:"quantity"`
		BuyPrice  string `clover:"buyPrice"`
		SellPrice string `clover:"sellPrice"`
	}{}

	doc := clover.NewDocument()
	doc.Set("exchange", "kucoin")
	doc.Set("status", "completed")
	doc.Set("quantity", "0.004")
	doc.Set("buyPrice", "94000")
	doc.Set("sellPrice", "98000")

	_, err := db.InsertOne("cycles", doc)
	if err != nil {
		return
	}

	docs, _ := db.Query("cycles").FindAll()
	for _, doc := range docs {
		err := doc.Unmarshal(cycle)
		if err != nil {
			return
		}
		fmt.Println(cycle.Id, cycle.Exchange, cycle.Status, cycle.Quantity, cycle.BuyPrice, cycle.SellPrice)
	}

	defer func(db *clover.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)
}
