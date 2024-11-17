// File: src/backend/internal/handlers/handlers.go

package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/db"
	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/models"
	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/utils"
)

var templates *template.Template

func GetTemplates() *template.Template {
	return templates
}

func SetTemplates(t *template.Template) {
	templates = t
}

func InitTemplates(pattern string) error {
	// Try multiple possible template locations
	patterns := []string{
		"frontend/templates/*.html",    // For production
		"../frontend/templates/*.html", // For local development
	}

	var templateErr error
	for _, p := range patterns {
		log.Printf("Trying template pattern: %s", p)
		t, err := template.ParseGlob(p)
		if err == nil {
			templates = t
			log.Printf("Successfully loaded templates from: %s", p)
			return nil
		}
		templateErr = err
	}

	// If we get here, no patterns worked
	log.Printf("Failed to load templates from any location. Last error: %v", templateErr)
	return templateErr
}

func init() {
	// This will be overridden in tests
	if err := InitTemplates(""); // pattern is ignored now
	err != nil {
		log.Printf("Warning: Failed to parse templates: %v", err)
	}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*sessions.Session)
	query := r.URL.Query().Get("q")
	var searchResults []map[string]interface{}
	var err error

	if query != "" {
		// Updated to use $1, $2 for PostgreSQL
		rows, err := db.DB.Query(
			"SELECT title, url, content FROM pages WHERE title LIKE $1 OR content LIKE $2",
			"%"+query+"%",
			"%"+query+"%",
		)
		if err != nil {
			log.Printf("Search query error: %v", err) // Added logging
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var title, url, content string
			err = rows.Scan(&title, &url, &content)
			if err != nil {
				log.Printf("Row scan error: %v", err) // Added logging
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			result := map[string]interface{}{
				"title":   title,
				"url":     url,
				"content": content,
			}
			searchResults = append(searchResults, result)
		}
	}

	data := map[string]interface{}{
		"Query":         query,
		"SearchResults": searchResults,
		"User":          session.Values["user"],
		"FlashMessages": session.Flashes(),
	}

	err = templates.ExecuteTemplate(w, "search", data)
	if err != nil {
		log.Printf("Template execution error: %v", err) // Added logging
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
	session.Save(r, w)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*sessions.Session)
	data := map[string]interface{}{
		"Flashes": session.Flashes(),
		"User":    session.Values["user"],
	}

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		user, err := db.GetUser(username)
		if err != nil {
			if err == db.ErrUserNotFound {
				data["Error"] = "Invalid username or password"
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		} else if user != nil && utils.CheckPasswordHash(password, user.Password) {
			session.Values["user"] = user
			session.Values["user_id"] = user.ID
			session.AddFlash("You were logged in")
			err = session.Save(r, w)
			if err != nil {
				log.Printf("Error saving session: %v", err)
				http.Error(w, "Error saving session", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			data["Error"] = "Invalid username or password"
		}
		data["Username"] = username
	}

	err := templates.ExecuteTemplate(w, "login", data)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*sessions.Session)
	data := map[string]interface{}{
		"Flashes": session.Flashes(),
		"User":    session.Values["user"],
	}

	if r.Method == "POST" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		password2 := r.FormValue("password2")

		if password != password2 {
			data["Error"] = "The two passwords do not match"
		} else {
			// First check if user exists
			user, err := db.GetUser(username)
			if err != nil && err != db.ErrUserNotFound {
				log.Printf("Database error checking user: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if user != nil {
				data["Error"] = "Username already taken"
			} else {
				err := db.CreateUser(username, email, password)
				if err != nil {
					// Check for unique constraint violation
					if strings.Contains(err.Error(), "unique constraint") {
						if strings.Contains(err.Error(), "users_username_key") {
							data["Error"] = "Username already taken"
						} else if strings.Contains(err.Error(), "users_email_key") {
							data["Error"] = "Email already registered"
						} else {
							data["Error"] = "Username or email already exists"
						}
					} else {
						log.Printf("Error creating user: %v", err)
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
						return
					}
				} else {
					session.AddFlash("You were successfully registered and can login now")
					session.Save(r, w)
					http.Redirect(w, r, "/login", http.StatusSeeOther)
					return
				}
			}
		}
		// Preserve form data on error
		data["Username"] = username
		data["Email"] = email
	}

	err := templates.ExecuteTemplate(w, "register", data)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
	session.Save(r, w)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*sessions.Session)

	// Clear all session values
	for key := range session.Values {
		delete(session.Values, key)
	}

	// Expire the cookie
	session.Options.MaxAge = -1

	log.Printf("Logging out user. Session values before save: %v", session.Values)

	err := session.Save(r, w)
	if err != nil {
		log.Printf("Error saving session during logout: %v", err)
		http.Error(w, "Error during logout", http.StatusInternalServerError)
		return
	}

	// Add flash message before redirect
	session.AddFlash("You have been successfully logged out")
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Error saving flash message: %v", err)
	}

	log.Printf("User successfully logged out, redirecting to home")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// AboutHandler renders the about template
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*sessions.Session)
	data := map[string]interface{}{
		"User": session.Values["user"], // Pass User to the template
	}

	err := templates.ExecuteTemplate(w, "about", data)
	if err != nil {
		log.Printf("Error rendering about template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*sessions.Session)

	weatherData, err := models.FetchWeather("Copenhagen")
	if err != nil {
		log.Printf("Error fetching weather data: %v", err)
		http.Error(w, "Error fetching weather data", http.StatusInternalServerError)
		return
	}

	templateData := map[string]interface{}{
		"Temperature": fmt.Sprintf("%.1f", weatherData.Main.Temp),
		"Condition":   weatherData.Weather[0].Main,
		"Description": weatherData.Weather[0].Description,
		"Location":    weatherData.Name,
		"User":        session.Values["user"],
	}

	err = templates.ExecuteTemplate(w, "weather", templateData)
	if err != nil {
		log.Printf("Error rendering weather template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
