// File src/backend/main.go:

package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/db"
	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/handlers"
	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/middleware"
)

var (
	store *sessions.CookieStore
)

func init() {
	// Load .env file first
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize session store
	sessionKey := "your-session-key"
	if sessionKey == "" {
		log.Fatal("SESSION_KEY environment variable is required")
	}
	store = sessions.NewCookieStore([]byte(sessionKey))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
	}

}

func main() {
	// Initialize the database
	if err := db.ConnectDB(); err != nil {
		log.Fatalf("Could not connect to lacal database: %v", err)
	}
	log.Println("Local database connected successfully.")

	// Create a new router
	r := mux.NewRouter()

	// Apply global middlewares
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.PrometheusMiddleware)
	r.Use(middleware.SessionMiddleware(store))

	// Set up routes
	r.HandleFunc("/", handlers.SearchHandler).Methods("GET")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("GET", "POST")
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("GET", "POST")
	r.HandleFunc("/about", handlers.AboutHandler).Methods("GET")
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("GET")
	r.HandleFunc("/weather", handlers.WeatherHandler).Methods("GET")
	r.HandleFunc("/greeting", handlers.Greeting).Methods("GET")

	// Prometheus metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/search", handlers.Search).Methods("GET")
	api.HandleFunc("/register", handlers.Register).Methods("POST")
	api.HandleFunc("/login", handlers.Login).Methods("POST")
	api.HandleFunc("/logout", handlers.Logout).Methods("GET")
	api.HandleFunc("/weather", handlers.Weather).Methods("GET")

	// Serve static files
	fs := http.FileServer(http.Dir("../frontend/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Start the server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
