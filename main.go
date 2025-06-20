package main

import (
	"fmt"
	"github.com/ostafen/clover"
	"main/commands"
	"main/database"
	"os"
	"slices"
)

func menu() {
	fmt.Println("")
	fmt.Println("Cryptomancien - BOT SPOT - v3.0.6 - beta")
	fmt.Println("")
	fmt.Println("--new			-n		Start new cycle")
	fmt.Println("--update		-u		Update running cycles")
	fmt.Println("--server		-s		Start local server")
	fmt.Println("--cancel		-c		Cancel cycle by id - Example: -c 123")
	fmt.Println("--auto			-a		Mode auto")
	fmt.Println("--clear 		-cl		Clear range (start end) - Example: -cl 12 36")
	fmt.Println("--export		-e		Export CSV file")
	fmt.Println("")
}

func initialize() error {
	commands.CreateConfigFileIfNotExists()
	commands.LoadDotEnv()

	db, err := database.InitDatabase()
	if err != nil {
		panic(err)
	}
	defer func(db *clover.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	return nil
}

func main() {
	err := initialize()
	if err != nil {
		panic("Error initializing database: %v")
	}

	args := os.Args[1:]

	actions := map[string]func(){
		"--new":    commands.New,
		"-n":       commands.New,
		"--update": commands.Update,
		"-u":       commands.Update,
		"--server": commands.Server,
		"-s":       commands.Server,
		"--cancel": commands.Cancel,
		"-c":       commands.Cancel,
		"--auto":   commands.Auto,
		"-a":       commands.Auto,
		"--clear":  commands.Clear,
		"-cl":      commands.Clear,
		"--export": func() { commands.Export(true) },
		"-e":       func() { commands.Export(true) },
	}

	for key, action := range actions {
		if slices.Contains(args, key) {
			action()
			return
		}
	}

	menu()
}
