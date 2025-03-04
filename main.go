package main

import (
	"fmt"
	"github.com/fatih/color"
	"main/commands/cancel"
	"main/commands/new"
	"main/commands/server"
	"main/commands/update"
	"main/database"
	"main/utils"
	"regexp"
)

func menu() {
	fmt.Println("")
	fmt.Println("Cryptomancien - BOT SPOT - v3.0.0 - beta")
	fmt.Println("")
	fmt.Println("--new			-n		Start new cycle")
	fmt.Println("--update		-u		Update running cycles")
	fmt.Println("--server		-s		Start local server")
	fmt.Println("--cancel		-c		Cancel cycle by id - Example: -c=123")
	//fmt.Println("--auto			-a		Mode auto")
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
		color.Magenta("Start new cycle")
		utils.CheckPremium()
		new.New()
	case "--update", "-u":
		color.Magenta("Updating running cycle...")
		//CheckPremium()
		update.Update()
	case regexp.MustCompile(`^--cancel=(\d+)$`).FindString(lastArg), regexp.MustCompile(`^-c=(\d+)$`).FindString(lastArg):
		cancel.Cancel()
	case "--server", "-s":
		server.Serve()
	default:
		menu()
	}
}
