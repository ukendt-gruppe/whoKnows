package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/db"
)

type SearchResponse struct {
	Data []map[string]interface{} `json:"data"`
}

type StandardResponse struct {
	Data map[string]interface{} `json:"data"`
}

type AuthResponse struct {
	StatusCode *int    `json:"statusCode,omitempty"`
	Message    *string `json:"message,omitempty"`
}

// APISearchHandler handles the /api/search endpoint
func APISearchHandler(w http.ResponseWriter, r *http.Request) {
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

	response := SearchResponse{Data: searchResults}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// APIWeatherHandler handles the /api/weather endpoint
func APIWeatherHandler(w http.ResponseWriter, r *http.Request) {
	// Mock weather data
	weatherData := map[string]interface{}{
		"temperature": 22,
		"condition":   "Sunny",
	}

	response := StandardResponse{Data: weatherData}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// APILoginHandler handles the /api/login endpoint
func APILoginHandler(w http.ResponseWriter, r *http.Request) {
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

	// Simple mock login without actual authentication
	if username != "" && password != "" {
		statusCode := http.StatusOK
		message := "Login successful"
		response := AuthResponse{StatusCode: &statusCode, Message: &message}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
	}
}

// APIRegisterHandler handles the /api/register endpoint
func APIRegisterHandler(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Simple mock registration without database interaction
	statusCode := http.StatusOK
	message := "User registered successfully"
	response := AuthResponse{StatusCode: &statusCode, Message: &message}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// APILogoutHandler handles the /api/logout endpoint
func APILogoutHandler(w http.ResponseWriter, r *http.Request) {
	statusCode := http.StatusOK
	message := "Logout successful"
	response := AuthResponse{StatusCode: &statusCode, Message: &message}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}