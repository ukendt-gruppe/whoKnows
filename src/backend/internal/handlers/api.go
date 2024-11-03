// File: src/backend/internal/handlers/api.go

package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/db"
	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/models"
	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/utils"
)

// SearchResponse represents the response for the search API
type SearchResponse struct {
	Data []map[string]interface{} `json:"data"`
}

// StandardResponse represents a standard API response
type StandardResponse struct {
	Data map[string]interface{} `json:"data"`
}

// AuthResponse represents the response for authentication-related APIs
type AuthResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// RequestValidationError represents validation errors
type RequestValidationError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// Search handles the /api/search endpoint
func Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	language := r.URL.Query().Get("language")
	if language == "" {
		language = "en"
	}

	var searchResults []map[string]interface{}
	var err error

	if query != "" {
		searchResults, err = db.QueryDB("SELECT * FROM pages WHERE language = ? AND content LIKE ?", language, "%"+query+"%")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	// Ensure empty data array if no results
	if searchResults == nil {
		searchResults = []map[string]interface{}{}
	}

	response := SearchResponse{Data: searchResults}
	utils.JSONResponse(w, http.StatusOK, response) // Ensure status 200
}

// Weather handles the /api/weather endpoint
func Weather(w http.ResponseWriter, r *http.Request) {
	weatherData, err := models.FetchWeather("Copenhagen")
	if err != nil {
		http.Error(w, "Error fetching weather data", http.StatusInternalServerError)
		return
	}

	response := utils.StandardResponse{
		Data: map[string]interface{}{
			"temperature": weatherData.Main.Temp,
			"condition":   weatherData.Weather[0].Main,
			"description": weatherData.Weather[0].Description,
			"location":    weatherData.Name,
		},
	}
	utils.JSONResponse(w, http.StatusOK, response)
}

// Register handles the /api/register endpoint
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		utils.JSONResponse(w, http.StatusUnprocessableEntity, RequestValidationError{
			StatusCode: http.StatusUnprocessableEntity, // Return 422 for validation error
			Message:    "All fields are required",
		})
		return
	}

	// Registration logic here

	response := AuthResponse{
		StatusCode: http.StatusOK,
		Message:    "User registered successfully",
	}
	utils.JSONResponse(w, http.StatusOK, response)
}

// Login handles the /api/login endpoint
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		utils.JSONResponse(w, http.StatusUnprocessableEntity, RequestValidationError{
			StatusCode: http.StatusUnprocessableEntity, // Return 422 for missing fields
			Message:    "Username and password are required",
		})
		return
	}

	// Implement authentication logic here!
	// If the credentials are correct, return the success response
	if username != "" && password != "" {
		response := AuthResponse{
			StatusCode: http.StatusOK, // Ensure status 200
			Message:    "Login successful",
		}
		utils.JSONResponse(w, http.StatusOK, response)
	} else {
		response := AuthResponse{
			StatusCode: http.StatusUnauthorized, // Return 401 for invalid login
			Message:    "Invalid username or password",
		}
		utils.JSONResponse(w, http.StatusUnauthorized, response)
	}
}

// Logout handles the /api/logout endpoint
func Logout(w http.ResponseWriter, r *http.Request) {
	response := AuthResponse{
		StatusCode: http.StatusOK,
		Message:    "Logout successful",
	}
	utils.JSONResponse(w, http.StatusOK, response)
}

func Greeting(w http.ResponseWriter, r *http.Request) {
	greeting := os.Getenv("ENV_GREETING")
	log.Printf("Current ENV_GREETING value: %q", greeting) // Debug line

	if greeting == "" {
		greeting = "No greeting set in environment"
	}
	w.Write([]byte(greeting))
}
