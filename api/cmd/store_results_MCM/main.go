package main

import (
	"api/db"
	"log"
)

func main() {
	dbc, err := db.New()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbc.Close()

	err = dbc.InitTablesMCM()
	if err != nil {
		log.Fatalf("Failed to initialize tables: %v", err)
	}

	// Import data into the 'actas_mcm' table
	err = dbc.ImportCSVForMariaCorina()
	if err != nil {
		log.Fatalf("Failed to import data: %v", err)
	}
}
