package commands

import (
	"github.com/joho/godotenv"
	"testing"
)

func TestList(t *testing.T) {
	err := godotenv.Load("../bot.conf")
	if err != nil {
		t.Fatalf("Error loading ../bot.conf: %v", err)
	}

	List()
}
