package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/db"
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
	StatusCode *int    `json:"statusCode,omitempty"`
	Message    *string `json:"message,omitempty"`
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)  // Ensure status 200
	json.NewEncoder(w).Encode(response)
}

// Weather handles the /api/weather endpoint
func Weather(w http.ResponseWriter, r *http.Request) {
	weatherData := map[string]interface{}{
		"temperature": 22,
		"condition":   "Sunny",
	}
	response := StandardResponse{Data: weatherData}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)  // Return 422 for validation error
		json.NewEncoder(w).Encode(RequestValidationError{
			StatusCode: 422,
			Message:    "All fields are required",
		})
		return
	}

	// Registration logic here

	statusCode := http.StatusOK
	message := "User registered successfully"
	response := AuthResponse{StatusCode: &statusCode, Message: &message}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)  // Ensure status 200
	json.NewEncoder(w).Encode(response)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)  // Return 422 for missing fields
		json.NewEncoder(w).Encode(RequestValidationError{
			StatusCode: 422,
			Message:    "Username and password are required",
		})
		return
	}

	// Implement authentication logic here
	// If the credentials are correct, return the success response
	if username != "" && password != "" {
		statusCode := http.StatusOK
		message := "Login successful"
		response := AuthResponse{StatusCode: &statusCode, Message: &message}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)  // Ensure status 200
		json.NewEncoder(w).Encode(response)
	} else {
		// Invalid credentials
		statusCode := http.StatusUnauthorized  // Assign constant to variable
		message := "Invalid username or password"  // Assign string to variable
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)  // Return 401 for invalid login
		json.NewEncoder(w).Encode(AuthResponse{
			StatusCode: &statusCode,
			Message:    &message,  // Use pointer to variable
		})
	}
}

// Logout handles the /api/logout endpoint
func Logout(w http.ResponseWriter, r *http.Request) {
	statusCode := http.StatusOK
	message := "Logout successful"
	response := AuthResponse{StatusCode: &statusCode, Message: &message}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}