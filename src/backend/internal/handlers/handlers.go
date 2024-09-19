package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/ukendt-gruppe/whoKnows/src/go-rewrite/internal/db"
)

var templates = template.Must(template.ParseGlob("../frontend/templates/*.html"))

// SearchHandler is responsible for processing HTTP requests related to search queries.
// It fetches the query parameter from the URL, interacts with the SQLite database, and renders the search results page.
// This function logs any errors encountered during the process, following best practices for observability.
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the search query from the request URL.
	query := r.URL.Query().Get("q")

	// Initialize the search results slice and error variable.
	var searchResults []map[string]interface{}
	var err error

	// If a query exists, proceed to search the database for relevant records.
	if query != "" {
		// Query the database for pages where the title or content matches the search term.
		// Using parameterized queries to prevent SQL injection, ensuring security.
		searchResults, err = db.QueryDB("SELECT title, url, content FROM pages WHERE title LIKE ? OR content LIKE ?", "%"+query+"%", "%"+query+"%")

		// Implementing proper error handling to ensure any database query failures are logged for future investigation (troubleshooting).
		if err != nil {
			log.Printf("Database query failed: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	// Prepare the data structure to pass to the template renderer.
	// This includes both the search query and the results fetched from the database.
	data := map[string]interface{}{
		"Query":         query,
		"SearchResults": searchResults,
	}

	// Render the "search" template, injecting the data (query and results) into the HTML.
	// Observability: Log any rendering errors to monitor potential issues with the templating system.
	err = templates.ExecuteTemplate(w, "search", data)
	if err != nil {
		log.Printf("Error rendering search template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

// LoginHandler renders the login template
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Render the "login" template without any additional data
	err := templates.ExecuteTemplate(w, "login", nil)
	if err != nil {
		log.Printf("Error rendering login template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

// RegisterHandler renders the register template
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Render the "register" template without any additional data
	err := templates.ExecuteTemplate(w, "register", nil)
	if err != nil {
		log.Printf("Error rendering register template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

// AboutHandler renders the about template
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	// Render the "about" template without any additional data
	err := templates.ExecuteTemplate(w, "about", nil)
	if err != nil {
		log.Printf("Error rendering about template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

// TestHandler renders the test template w/o errors
func TestHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "test", nil)
}
