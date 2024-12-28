package main

import (
	"log"

	"github.com/jaxxk/anki-cards-generator/cmd"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cmd.Execute()
}
