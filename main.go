package main

import (
	"fmt"
	"main/commands"
	"main/database"
	"regexp"
)

func menu() {
	fmt.Println("")
	fmt.Println("Cryptomancien - BOT SPOT - v3.0.2 - beta")
	fmt.Println("")
	fmt.Println("--new			-n		Start new cycle")
	fmt.Println("--update		-u		Update running cycles")
	fmt.Println("--server		-s		Start local server")
	fmt.Println("--cancel		-c		Cancel cycle by id - Example: -c=123")
	fmt.Println("--auto			-a		Mode auto")
	fmt.Println("")
}

func initialize() {
	commands.CreateConfigFileIfNotExists()
	commands.LoadDotEnv()
	database.InitDatabase()
}

func main() {
	initialize()

	lastArg := commands.GetLastArg()
	switch lastArg {

	case "--new", "-n":
		commands.New()
	case "--update", "-u":
		commands.Update()
	case regexp.MustCompile(`^--cancel=(\d+)$`).FindString(lastArg), regexp.MustCompile(`^-c=(\d+)$`).FindString(lastArg):
		commands.Cancel()
	case "--auto", "-a":
		commands.Auto()
	case "--server", "-s":
		commands.Server()
	default:
		menu()
	}
}
