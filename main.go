package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/rekib0023/event-horizon-gateway/controller"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	controller.Start()

}
