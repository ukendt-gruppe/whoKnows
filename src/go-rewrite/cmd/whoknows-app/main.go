package main

import (
	"log"
    "github.com/ukendt-gruppe/whoKnows/internal/db"
)

func main() {
	if err := db.InitDB("schema.sql"); err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	log.Println("Database initialized successfully.")
}