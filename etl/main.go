package main

import (
	"context"
	"log"

	"github.com/thalesmacedo1/covid-etl/config"
	"github.com/thalesmacedo1/covid-etl/services"
)

func main() {
	// Load configuration
	err := config.LoadConfig(".env.example")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create ETL service
	etlService, err := services.NewETLService()
	if err != nil {
		log.Fatalf("Failed to initialize ETL service: %v", err)
	}

	// Run ETL process
	ctx := context.Background()
	if err := etlService.Run(ctx); err != nil {
		log.Fatalf("ETL process failed: %v", err)
	}

	log.Println("ETL process completed successfully.")
}
