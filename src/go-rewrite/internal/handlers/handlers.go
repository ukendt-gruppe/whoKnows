package handlers

import (
    "html/template"
    "net/http"
    "log"
)

var templates = template.Must(template.ParseGlob("src/go-rewrite/frontend/templates/*.html"))

// TestHandler renders the test template
func TestHandler(w http.ResponseWriter, r *http.Request) {
    // Render the "test" template without any additional data
    err := templates.ExecuteTemplate(w, "test", nil)
    if err != nil {
        log.Printf("Error rendering test template: %v", err)
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