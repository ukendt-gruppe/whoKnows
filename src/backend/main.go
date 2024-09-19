package main

import (
	"log"
	"net/http"

	"github.com/ukendt-gruppe/whoKnows/src/go-rewrite/internal/db"
	"github.com/ukendt-gruppe/whoKnows/src/go-rewrite/internal/handlers"
)

func main() {
    if err := db.InitDB("./schema.sql"); err != nil {
        log.Fatalf("Could not initialize database: %v", err)
    }
    log.Println("Database initialized successfully.")

    // Set up routes
    http.HandleFunc("/", handlers.SearchHandler)
    http.HandleFunc("/login", handlers.LoginHandler)
    http.HandleFunc("/register", handlers.RegisterHandler)
    http.HandleFunc("/about", handlers.AboutHandler)
    http.HandleFunc("/test", handlers.TestHandler)


    // Serve static files
    fs := http.FileServer(http.Dir("../frontend/static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    // Start the server
    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}