package main

import (
	"fmt"
	"main/commands"
	"main/database"
	"main/server"
	"main/utils"
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
	utils.CreateConfigFileIfNotExists()
	utils.LoadDotEnv()
	database.InitDatabase()
}

func main() {
	initialize()

	lastArg := utils.GetLastArg()
	switch lastArg {

	case "--new", "-n":
		fmt.Println("Start new cycle")
		utils.CheckPremium()
		commands.New()
	case "--update", "-u":
		fmt.Println("Updating running cycle...")
		//CheckPremium()
		commands.Update()
	case "--server", "-s":
		server.Serve()
	default:
		menu()
	}
}
