package main

import (
	"log"

	sharedBootstrap "github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/bootstrap"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/database"
)

func main() {
	depedency, err := sharedBootstrap.NewDeps()
	if err != nil {
		log.Fatalf("error injecting depedency %v", err)
	}

	// Migrate and seed the database
	if err := database.Migrate(depedency.DB); err != nil {
		log.Fatalf("error migrating database: %v", err)
	}

	if err := database.Seeding(depedency.DB); err != nil {
		log.Fatalf("error seeding database: %v", err)
	}

	log.Println("success migrating and seeding database")
}
