package middleware

import (
	"context"
	"net/http"
)

// MockAuthMiddleware simulates authentication without actually checking credentials
func MockAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Simulate a logged-in user
		ctx := context.WithValue(r.Context(), "user", map[string]string{"username": "mock_user"})
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// LoggingMiddleware logs incoming requests
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the request here
		// For example: log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	}
}
