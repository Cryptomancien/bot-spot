package commands

import (
	"encoding/csv"
	"fmt"
	"github.com/fatih/color"
	"main/database"
	"os"
	"path/filepath"
	"time"
)

func Export() {
	// Format current timestamp like Ruby (e.g. "2025-05-13 14:30:00")
	timestamp := time.Now().Format("2006-01-02 15-04-05")

	// Create "exports" folder if it doesn't exist
	if err := os.MkdirAll("exports", os.ModePerm); err != nil {
		panic(fmt.Errorf("failed to create exports folder: %w", err))
	}

	// Create CSV file in the "exports" folder
	fileName := fmt.Sprintf("exports/%s.csv", timestamp)

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
		"BuyPrice", "SellPrice", "Gain", "BuyId", "SellId", "_id",
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
			fmt.Sprintf("%v", cycle.Get("gain")),
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
