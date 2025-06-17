package commands

import (
	"encoding/csv"
	"fmt"
	"github.com/fatih/color"
	"github.com/ostafen/clover"
	"log"
	"main/database"
	"os"
	"path/filepath"
	"time"
)

func fileNamePrefix() string {
	// Create an "exports" folder if it doesn't exist
	if err := os.MkdirAll("exports", os.ModePerm); err != nil {
		panic(fmt.Errorf("failed to create exports folder: %w", err))
	}

	timestamp := time.Now().Format("2006-01-02 15-04-05")

	filePrefix := fmt.Sprintf("exports/%s", timestamp)
	return filePrefix
}

func ToCSV() {
	fileName := fileNamePrefix() + ".csv"

	color.Yellow("Export data to CSV file: " + fileName)

	file, err := os.Create(filepath.Clean(fileName))
	if err != nil {
		panic(fmt.Errorf("failed to create file: %w", err))
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(fmt.Errorf("failed to close file: %w", err))
		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"IdInt", "Exchange", "Status", "Quantity",
		"BuyPrice", "SellPrice", "Gain USD", "BuyId", "SellId", "_id",
	}
	if err := writer.Write(header); err != nil {
		panic(fmt.Errorf("failed to write header: %w", err))
	}

	// Write each row
	cycles := database.List()
	for _, cycle := range cycles {
		row := []string{
			fmt.Sprintf("%v", cycle.Get("idInt")),
			fmt.Sprintf("%v", cycle.Get("exchange")),
			fmt.Sprintf("%v", cycle.Get("status")),
			fmt.Sprintf("%v", cycle.Get("quantity")),
			fmt.Sprintf("%v", cycle.Get("buyPrice")),
			fmt.Sprintf("%v", cycle.Get("sellPrice")),

			fmt.Sprintf("%v", CalcAbsoluteGainByCycle(cycle)),

			fmt.Sprintf("%v", cycle.Get("buyId")),
			fmt.Sprintf("%v", cycle.Get("sellId")),
			fmt.Sprintf("%v", cycle.Get("_id")),
		}

		if err := writer.Write(row); err != nil {
			panic(fmt.Errorf("failed to write row: %w", err))
		}
	}
	color.Green("Successfully Export data to CSV file: " + fileName)
}

func ToJSON() {
	fileName := fileNamePrefix() + ".json"

	db := database.GetDB()
	defer func(db *clover.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("can't close database %s", err)
		}
	}(db)

	err := db.ExportCollection(database.CollectionName, fileName)
	if err != nil {
		log.Fatal("Can't export collection to JSON file: " + fileName)
	}
	color.Green("Successfully Export data to JSON file: " + fileName)
}

func Export() {
	ToCSV()
	ToJSON()
}
