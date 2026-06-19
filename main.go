package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/zwinslett/speed-daemon/cmd"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load(home + "/.config/speed-daemon/.env")
	if err != nil {
		log.Fatal("Error loading .env file.")
	}
	cmd.Execute()
}
