package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/db"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func SessionMiddleware(store sessions.Store) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            session, _ := store.Get(r, "session-name")
            ctx := context.WithValue(r.Context(), "session", session)
            
            if userID, ok := session.Values["user_id"].(int); ok {
                user, err := db.GetUser(userID)
                if err == nil && user != nil {
                    session.Values["user"] = user
                } else {
                    delete(session.Values, "user")
                    delete(session.Values, "user_id")
                }
            }
            
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func PasswordResetCheckMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Skip check for certain paths
        if r.URL.Path == "/reset-password" || r.URL.Path == "/login" || r.URL.Path == "/logout" || r.URL.Path == "/static/" {
            next.ServeHTTP(w, r)
            return
        }

        session := r.Context().Value("session").(*sessions.Session)
        if user, ok := session.Values["user"].(*db.User); ok && user != nil {
            if user.NeedsPasswordReset {
                http.Redirect(w, r, "/reset-password", http.StatusSeeOther)
                return
            }
        }
        
        next.ServeHTTP(w, r)
    })
}
