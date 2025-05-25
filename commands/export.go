package commands

import (
	"encoding/csv"
	"fmt"
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

	// Get cycles from database
	cycles := database.List()
	if len(cycles) == 0 {
		return // nothing to write
	}

}
