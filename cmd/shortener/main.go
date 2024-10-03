package main

import (
	"log"
	"os"

	"github.com/orekhovskiy/shrtn/config"
	"github.com/orekhovskiy/shrtn/internal/app"
)

func main() {
	conf, err := config.InitializeConfig()
	if err != nil {
		log.Printf("error initializing config: %v", err)
		os.Exit(1)
	}

	app.Run(conf)

	if err != nil {
		log.Printf("error starting server: %v", err)
		os.Exit(1)
	}
}
