package commands

import (
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"log"
	"main/database"
	"os"
)

func Restore() {
	color.Cyan("Restoring database from export file")
	err := godotenv.Load("../bot.conf")
	if err != nil {
		panic("Error loading config file")
	}

	exportFile := os.Getenv("EXPORT_FILE")
	if exportFile == "" {
		color.Red("Missing environment variable: EXPORT_FILE")
		return
	}

	path := "../exports/" + exportFile

	color.Yellow("Importing data from file: " + path)

	db := database.GetDB()

	err = db.DropCollection(database.CollectionName)
	if err != nil {
		log.Println(err)
		return
	}

	err = db.ImportCollection(database.CollectionName, path)
	if err != nil {
		log.Println(err)
		color.Red("Error importing collection: " + path)
		return
	}

	color.Green("Database restored")
}
