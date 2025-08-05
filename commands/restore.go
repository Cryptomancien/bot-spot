package commands

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"log"
	"main/database"
	"os"
	"path/filepath"
)

func Restore() {

	// Load config file
	err := godotenv.Load("bot.conf")
	if err != nil {
		log.Fatalf("Error loading ../bot.conf: %v", err)
	}

	color.Cyan("Restoring database from export file, Follow this instructions:")
	fmt.Println("1) Find latest export file in exports folder, something like this: 2025-08-05 11-16-33.json")
	fmt.Println("2) Add EXPORT_FILE to config file, something like this: EXPORT_FILE=2025-08-05 11-16-33.json")
	fmt.Println("")

	// Checking requirement
	exportFile := os.Getenv("EXPORT_FILE")
	if exportFile == "" {
		color.Red("Missing environment variable: EXPORT_FILE")
		return
	}
	color.Green("Export file found in bot.conf: " + exportFile)

	// Find the export file
	workDir, err := os.Getwd()
	pathJSON := filepath.Join(workDir, "exports", exportFile)
	//fmt.Println("path:", exportDatabaseJSON)

	// Ensure file is present
	if _, err := os.Stat(pathJSON); errors.Is(err, os.ErrNotExist) {
		color.Red("Export file not found: " + pathJSON)
		return
	}
	color.Green("Export file found: " + pathJSON)

	// Find the DB path
	homePath, _ := os.UserHomeDir()
	dbFolder := "/cryptomancien/bot-v3"
	joinedFolder := filepath.Join(homePath, dbFolder)

	// Destroy DB folder
	err = os.RemoveAll(joinedFolder)
	if err != nil {
		log.Println(err)
		return
	} else {
		color.Green("DB folder successfully removed")
	}

	// Init db
	_, err = database.InitDatabase()
	if err != nil {
		panic(err)
	}

	db := database.GetDB()
	err = db.DropCollection(database.CollectionName)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	err = db.ImportCollection(database.CollectionName, pathJSON)
	if err != nil {
		log.Println(err)
		color.Red("Error importing collection: " + pathJSON)
		return
	}
	color.Green("Database successfully restored from export file: " + pathJSON)

	fmt.Println("\nrun command: go run . -s")
}
