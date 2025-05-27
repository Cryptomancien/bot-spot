package main

import (
	"fmt"
	"main/commands"
	"main/database"
	"os"
	"slices"
)

func menu() {
	fmt.Println("")
	fmt.Println("Cryptomancien - BOT SPOT - v3.0.4 - beta")
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

func initialize() {
	commands.CreateConfigFileIfNotExists()
	commands.LoadDotEnv()
	database.InitDatabase()
}

func main() {
	initialize()

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
		"--export": commands.Export,
		"-e":       commands.Export,
	}

	for key, action := range actions {
		if slices.Contains(args, key) {
			action()
			return
		}
	}

	menu()
}
