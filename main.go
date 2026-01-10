package main

import (
	"github.com/TwiLightDM/diploma-course-service/internal/app"
	"github.com/TwiLightDM/diploma-course-service/internal/config"
	"log"
)

func main() {
	cfg := config.Load()

	if err := app.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
