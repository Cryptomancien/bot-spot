package main

import (
	"fmt"
	"github.com/fatih/color"
	"main/commands"
	"main/database"
	"main/server"
)

func menu() {
	fmt.Println("")
	fmt.Println("Cryptomancien - BOT SPOT - v3.0.0 - beta")
	fmt.Println("")
	fmt.Println("--new			-n		Start new cycle")
	fmt.Println("--update		-u		Update running cycles")
	fmt.Println("--server		-s		Start local server")
	fmt.Println("--cancel		-cc		Cancel cycle by id")
	fmt.Println("")
}

func initialize() {
	CreateConfigFileIfNotExists()
	LoadDotEnv()
	database.InitDatabase()
}

func main() {
	initialize()

	lastArg := GetLastArg()
	switch lastArg {

	case "--new", "-n":
		color.Green("Start new cycle")
		fmt.Println("")
		CheckPremium()
		commands.New()
	case "--update", "-u":
		color.Green("Update running cycles")
		fmt.Println("")
		//CheckPremium()
		commands.Update()
	case "--server", "-s":
		server.Serve()
	default:
		menu()
	}
}
