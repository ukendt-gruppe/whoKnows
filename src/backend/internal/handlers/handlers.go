// File: src/backend/internal/handlers/handlers.go

package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

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
	var err error
	templates, err = template.ParseGlob(pattern)
	return err
}

func init() {
	// This will be overridden in tests
	err := InitTemplates("../frontend/templates/*.html")
	if err != nil {
		log.Printf("Warning: Failed to parse templates: %v", err)
	}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*sessions.Session)
	query := r.URL.Query().Get("q")
	var searchResults []map[string]interface{}
	var err error

	if query != "" {
		rows, err := db.DB.Query("SELECT title, url, content FROM pages WHERE title LIKE ? OR content LIKE ?", "%"+query+"%", "%"+query+"%")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var title, url, content string
			err = rows.Scan(&title, &url, &content)
			if err != nil {
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
			// Check username
			userByUsername, err := db.GetUser(username)
			if err != nil && err != db.ErrUserNotFound {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Check email
			userByEmail, err := db.GetUserByEmail(email)
			if err != nil && err != db.ErrUserNotFound {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if userByUsername != nil {
				data["Error"] = "The username is already taken"
			} else if userByEmail != nil {
				data["Error"] = "The email address is already registered"
			} else {
				err := db.CreateUser(username, email, password)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				session.AddFlash("You were successfully registered and can login now")
				session.Save(r, w)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		}
		data["Username"] = username
		data["Email"] = email
	}

	err := templates.ExecuteTemplate(w, "register", data)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
	session.Save(r, w)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*sessions.Session)

	// Clear all session values
	delete(session.Values, "user")
	delete(session.Values, "user_id")

	// Optional: Clear all session data
	session.Options.MaxAge = -1 // This will tell the browser to remove the cookie

	log.Printf("Logging out user. Session values after cleanup: %+v", session.Values)

	session.AddFlash("You were logged out")
	err := session.Save(r, w)
	if err != nil {
		log.Printf("Error saving session during logout: %v", err)
		http.Error(w, "Error saving session", http.StatusInternalServerError)
		return
	}

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
