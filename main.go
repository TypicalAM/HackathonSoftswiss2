package main

import (
	"log"

	"github.com/TypicalAM/HackathonSoftswiss2/config"
	"github.com/TypicalAM/HackathonSoftswiss2/models"
	"github.com/TypicalAM/HackathonSoftswiss2/routes"
)

func main() {
	// Read the config file
	cfg := config.New()

	// Connect to the database
	db, err := models.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Migrate the database
	err = models.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	// Set up the router
	router, err := routes.New(db, cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Run the app
	if err = router.Run(cfg.ListenPort); err != nil {
		log.Fatal(err)
	}
}
