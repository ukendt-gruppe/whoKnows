package handlers

import (
	"net/http"

	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/db"
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

// Search godoc
// @Summary Search for content
// @Description Perform a search in the database using a query and language.
// @Tags Search
// @Accept  json
// @Produce  json
// @Param q query string true "Search query"
// @Param language query string false "Language (default: en)"
// @Success 200 {object} SearchResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/search [get]
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

// Weather godoc
// @Summary Get weather data
// @Description Returns the current weather conditions.
// @Tags Weather
// @Produce  json
// @Success 200 {object} StandardResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/weather [get]
func Weather(w http.ResponseWriter, r *http.Request) {
	weatherData := map[string]interface{}{
		"temperature": 14,
		"condition":   "Rainy",
	}
	response := StandardResponse{Data: weatherData}
	utils.JSONResponse(w, http.StatusOK, response)
}

// Register handles the /api/register endpoint

// Register godoc
// @Summary Register a new user
// @Description Registers a new user by taking a username and password.
// @Tags Authentication
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {object} AuthResponse
// @Failure 422 {object} RequestValidationError
// @Failure 405 {object} ErrorResponse
// @Router /api/register [post]
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

// Login godoc
// @Summary Log in a user
// @Description Logs in a user with a username and password.
// @Tags Authentication
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {object} AuthResponse
// @Failure 401 {object} AuthResponse
// @Failure 422 {object} RequestValidationError
// @Failure 405 {object} ErrorResponse
// @Router /api/login [post]
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

	// Implement authentication logic here
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

// Logout godoc
// @Summary Log out a user
// @Description Logs out the current user.
// @Tags Authentication
// @Produce  json
// @Success 200 {object} AuthResponse
// @Router /api/logout [get]
func Logout(w http.ResponseWriter, r *http.Request) {
	response := AuthResponse{
		StatusCode: http.StatusOK,
		Message:    "Logout successful",
	}
	utils.JSONResponse(w, http.StatusOK, response)
}
