package handlers

import (
    "html/template"
    "net/http"
    "log"
)

var templates = template.Must(template.ParseGlob("src/go-rewrite/frontend/templates/*.html"))

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// Render the "search" template without any additional data
	err := templates.ExecuteTemplate(w, "search", nil)
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



