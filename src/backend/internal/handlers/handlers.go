package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/db"
	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/utils"
)

var templates = template.Must(template.ParseGlob("./frontend/templates/*.html"))

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
			user, err := db.GetUser(username)
			if err != nil && err != db.ErrUserNotFound {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if user != nil {
				data["Error"] = "The username is already taken"
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
    delete(session.Values, "user")
    log.Println("User logged out:", session.Values["user"])
    session.AddFlash("You were logged out")
    err := session.Save(r, w)
    if err != nil {
        http.Error(w, "Error saving session", http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// AboutHandler renders the about template
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*sessions.Session)
	data := map[string]interface{}{
		"User": session.Values["user"],  // Pass User to the template
	}

	err := templates.ExecuteTemplate(w, "about", data)
	if err != nil {
		log.Printf("Error rendering about template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*sessions.Session)
	weatherData := map[string]interface{}{
		"Temperature": 15,
		"Condition":   "Rainy",
		"Location":    "Copenhagen",
		"User":        session.Values["user"],  // Pass User to the template
	}

	err := templates.ExecuteTemplate(w, "weather", weatherData)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
