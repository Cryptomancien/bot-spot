package tools

import (
	"encoding/csv"
	"fmt"
	"log"
	"main/database"
	"os"
	"strconv"
)

func SaveExportToDatabase() {
	csvFile := os.Getenv("CSV_FILE")

	f, err := os.Open("../exports/" + csvFile)
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for index, record := range records {
		if index == 0 {
			continue
		}

		idInt, _ := strconv.Atoi(record[0])
		exchange := record[1]
		status := record[2]
		quantity, _ := strconv.ParseFloat(record[3], 64)
		buyPrice, _ := strconv.ParseFloat(record[4], 64)
		sellPrice, _ := strconv.ParseFloat(record[5], 64)
		buyId := record[7]
		sellId := record[8]

		cycleExists := database.GetByIdInt(idInt)

		if cycleExists != nil {
			fmt.Printf("Cycle %d already exists, skipping\n", idInt)
			continue
		}

		cycle := database.Cycle{
			IdInt:     int32(idInt),
			Exchange:  exchange,
			Status:    status,
			Quantity:  quantity,
			BuyPrice:  buyPrice,
			SellPrice: sellPrice,
			BuyId:     buyId,
			SellId:    sellId,
		}

		database.NewCycle(&cycle)
		fmt.Printf("Cycle %d successfully inserted in database\n", idInt)
	}
}
