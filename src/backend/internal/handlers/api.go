// File: src/backend/internal/handlers/api.go

package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/db"
	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/models"
	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/utils"
)

// SearchResponse represents the response for the search APIss
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
// @Summary Search content
// @Description Perform a search across pages and wiki_articles.
// @Tags search
// @Accept  json
// @Produce  json
// @Param q query string true "Search query"
// @Param language query string false "Language for search (default: en)"
// @Success 200 {object} SearchResponse
// @Failure 500 {object} RequestValidationError
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
		// Modified query to only check language in pages table
		searchResults, err = db.QueryDB(`
			SELECT title, url, content, 'page' as source 
			FROM pages 
			WHERE language = $1 AND (title LIKE $2 OR content LIKE $2)
			UNION ALL
			SELECT title, url, content, 'wiki' as source 
			FROM wiki_articles 
			WHERE title LIKE $2 OR content LIKE $2
			ORDER BY title
		`, language, "%"+query+"%")

		if err != nil {
			log.Printf("Search query error: %v", err)
			http.Error(w, intErr, http.StatusInternalServerError)
			return
		}
	}

	if searchResults == nil {
		searchResults = []map[string]interface{}{}
	}

	response := SearchResponse{Data: searchResults}
	utils.JSONResponse(w, http.StatusOK, response)
}

// Weather handles the /api/weather endpoint

// Weather godoc
// @Summary Get weather data
// @Description Fetch weather data for a specific location.
// @Tags weather
// @Produce  json
// @Success 200 {object} StandardResponse
// @Failure 500 {object} RequestValidationError
// @Router /api/weather [get]

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

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email, and password.
// @Tags auth
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Param username formData string true "Username"
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Success 201 {object} AuthResponse
// @Failure 422 {object} RequestValidationError
// @Failure 409 {object} AuthResponse
// @Failure 500 {object} RequestValidationError
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
	email := r.FormValue("email")
	password := r.FormValue("password")

	if username == "" || email == "" || password == "" {
		utils.JSONResponse(w, http.StatusUnprocessableEntity, RequestValidationError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "All fields are required",
		})
		return
	}

	// Check if user exists
	existingUser, err := db.GetUser(username)
	if err != nil && err != db.ErrUserNotFound {
		log.Printf("Database error checking user: %v", err)
		http.Error(w, intErr, http.StatusInternalServerError)
		return
	}

	if existingUser != nil {
		utils.JSONResponse(w, http.StatusConflict, AuthResponse{
			StatusCode: http.StatusConflict,
			Message:    "Username already exists",
		})
		return
	}

	// Create new user
	err = db.CreateUser(username, email, password)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		if err.Error() == "username already exists" {
			utils.JSONResponse(w, http.StatusConflict, AuthResponse{
				StatusCode: http.StatusConflict,
				Message:    "Username already exists",
			})
			return
		}
		http.Error(w, intErr, http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		StatusCode: http.StatusCreated,
		Message:    "User registered successfully",
	}
	utils.JSONResponse(w, http.StatusCreated, response)
}

// Login handles the /api/login endpoint

// Login godoc
// @Summary Log in a user
// @Description Log in a user with username and password.
// @Tags auth
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {object} AuthResponse
// @Failure 401 {object} AuthResponse
// @Failure 422 {object} RequestValidationError
// @Failure 500 {object} RequestValidationError
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
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Username and password are required",
		})
		return
	}

	user, err := db.GetUser(username)
	if err != nil {
		if err == db.ErrUserNotFound {
			utils.JSONResponse(w, http.StatusUnauthorized, AuthResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    "Invalid username or password",
			})
			return
		}
		log.Printf("Database error during login: %v", err)
		http.Error(w, intErr, http.StatusInternalServerError)
		return
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		utils.JSONResponse(w, http.StatusUnauthorized, AuthResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid username or password",
		})
		return
	}

	// Set session
	session := r.Context().Value("session").(*sessions.Session)
	session.Values["user"] = user
	session.Values["user_id"] = user.ID
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Error saving session: %v", err)
		http.Error(w, intErr, http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		StatusCode: http.StatusOK,
		Message:    "Login successful",
	}
	utils.JSONResponse(w, http.StatusOK, response)
}

// Logout handles the /api/logout endpoint

// Logout godoc
// @Summary Log out a user
// @Description Log out the current user and clear the session.
// @Tags auth
// @Produce  json
// @Success 200 {object} AuthResponse
// @Failure 500 {object} RequestValidationError
// @Router /api/logout [get]
func Logout(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*sessions.Session)

	// Clear all session values
	for key := range session.Values {
		delete(session.Values, key)
	}

	// Expire the cookie
	session.Options.MaxAge = -1

	log.Printf("API: Logging out user. Session values before save: %v", session.Values)

	err := session.Save(r, w)
	if err != nil {
		log.Printf("API: Error saving session during logout: %v", err)
		http.Error(w, "Error during logout", http.StatusInternalServerError)
		return
	}

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
