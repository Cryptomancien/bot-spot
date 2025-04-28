package commands

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"time"
)

func dotenvToDuration(key string) time.Duration {
	str := os.Getenv(key)
	if str == "" {
		log.Fatal("Missing environment variable: " + key)
	}

	// If the string has no unit, assume minutes
	if str[len(str)-1] < 'a' || str[len(str)-1] > 'z' {
		str += "m"
	}

	duration, err := time.ParseDuration(str)
	if err != nil {
		log.Fatal(err)
	}

	return duration
}

func startNewCycle() {
	duration := dotenvToDuration("AUTO_INTERVAL_MINUTES_NEW")

	color.Blue("Starting new cycle every %d minutes", int(duration.Minutes()))

	for range time.Tick(time.Minute * duration) {
		fmt.Println(time.Now().Format(time.RubyDate))
		New()
	}
}

func updateRunningCycles() {
	duration := dotenvToDuration("AUTO_INTERVAL_MINUTES_UPDATE")

	color.Blue("Updating running cycles every %d minutes", int(duration.Minutes()))

	for range time.Tick(time.Minute * duration) {
		fmt.Println(time.Now().Format(time.RubyDate))
		Update()
	}
}

func Auto() {
	color.Yellow("Starting Auto Mode - CTRL + C to exit")

	go startNewCycle()
	go updateRunningCycles()

	// Exit after 24 hours
	time.Sleep(time.Hour * 24)
}
