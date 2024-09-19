package main

import (
	"log"
	"net/http"

	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/db"
	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/handlers"
    "github.com/ukendt-gruppe/whoKnows/src/backend/internal/middleware"
)

func main() {
    if err := db.InitDB("./internal/db/schema.sql"); err != nil {
        log.Fatalf("Could not initialize database: %v", err)
    }
    log.Println("Database initialized successfully.")

    // Set up routes
    http.HandleFunc("/", middleware.LoggingMiddleware(handlers.SearchHandler))
    http.HandleFunc("/login", middleware.LoggingMiddleware(handlers.LoginHandler))
    http.HandleFunc("/register", middleware.LoggingMiddleware(handlers.RegisterHandler))
    http.HandleFunc("/about", middleware.LoggingMiddleware(handlers.AboutHandler))
    http.HandleFunc("/weather", middleware.LoggingMiddleware(handlers.TestHandler))

    // API routes
    http.HandleFunc("/api/search", middleware.LoggingMiddleware(handlers.APISearchHandler))
    http.HandleFunc("/api/weather", middleware.LoggingMiddleware(handlers.APIWeatherHandler))
    http.HandleFunc("/api/register", middleware.LoggingMiddleware(handlers.APIRegisterHandler))
    http.HandleFunc("/api/login", middleware.LoggingMiddleware(handlers.APILoginHandler))
    http.HandleFunc("/api/logout", middleware.LoggingMiddleware(handlers.APILogoutHandler))

    // Serve static files
    fs := http.FileServer(http.Dir("../frontend/static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    // Start the server
    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
