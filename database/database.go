package database

import (
	"errors"
	"github.com/ostafen/clover"
	"log"
	"os"
	"path/filepath"
)

func InitDatabase() {
	const collectionName = "cycles"
	databasePath := GetDatabasePath()
	db, _ := clover.Open(databasePath)
	collectionAlreadyExists, _ := db.HasCollection(collectionName)
	if !collectionAlreadyExists {
		err := db.CreateCollection(collectionName)
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
