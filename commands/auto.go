package commands

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"sync"
	"time"
)

func dotenvToDuration(key string) time.Duration {
	str := os.Getenv(key)
	if str == "" {
		log.Fatal("Missing environment variable: " + key)
	}
	if str[len(str)-1] < 'a' || str[len(str)-1] > 'z' {
		str += "m"
	}
	duration, err := time.ParseDuration(str)
	if err != nil {
		log.Fatal(err)
	}
	return duration
}

func startNewCycle(wg *sync.WaitGroup) {
	defer wg.Done()
	duration := dotenvToDuration("AUTO_INTERVAL_NEW")
	color.Magenta("Starting new cycle every %s", duration.String())
	for range time.Tick(duration) {
		fmt.Println(time.Now().Format(time.RubyDate))
		New()
	}
}

func updateRunningCycles(wg *sync.WaitGroup) {
	defer wg.Done()
	duration := dotenvToDuration("AUTO_INTERVAL_UPDATE")
	color.Magenta("Updating running cycles every %s", duration.String())
	for range time.Tick(duration) {
		fmt.Println(time.Now().Format(time.RubyDate))
		Update()
	}
}

func Auto() {
	color.Yellow("Starting Auto Mode - CTRL + C to exit")

	var wg sync.WaitGroup
	wg.Add(2)

	go startNewCycle(&wg)
	go updateRunningCycles(&wg)

	wg.Wait()
}
