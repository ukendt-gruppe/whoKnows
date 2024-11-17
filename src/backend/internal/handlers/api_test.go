// File: src/backend/internal/handlers/api_test.go

package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	SetupTestTemplates(t)
	tests := []struct {
		name           string
		method         string
		username       string
		password       string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid registration",
			method:         http.MethodPost,
			username:       "testuser",
			password:       "testpass",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"statusCode":200,"message":"User registered successfully"}`,
		},
		{
			name:           "Missing username",
			method:         http.MethodPost,
			username:       "",
			password:       "testpass",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"statusCode":422,"message":"All fields are required"}`,
		},
		{
			name:           "Missing password",
			method:         http.MethodPost,
			username:       "testuser",
			password:       "",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"statusCode":422,"message":"All fields are required"}`,
		},
		{
			name:           "Invalid method",
			method:         http.MethodGet,
			username:       "testuser",
			password:       "testpass",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method Not Allowed", // Note the \n at the end
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("username", tt.username)
			form.Add("password", tt.password)

			req, err := http.NewRequest(tt.method, "/api/register", strings.NewReader(form.Encode()))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(Register)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if strings.TrimSpace(rr.Body.String()) != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	SetupTestTemplates(t)
	tests := []struct {
		name           string
		method         string
		username       string
		password       string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid login",
			method:         http.MethodPost,
			username:       "testuser",
			password:       "testpass",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"statusCode":200,"message":"Login successful"}`,
		},
		{
			name:           "Missing username",
			method:         http.MethodPost,
			username:       "",
			password:       "testpass",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"statusCode":422,"message":"Username and password are required"}`,
		},
		{
			name:           "Missing password",
			method:         http.MethodPost,
			username:       "testuser",
			password:       "",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"statusCode":422,"message":"Username and password are required"}`,
		},
		{
			name:           "Invalid method",
			method:         http.MethodGet,
			username:       "testuser",
			password:       "testpass",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method Not Allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("username", tt.username)
			form.Add("password", tt.password)

			req, err := http.NewRequest(tt.method, "/api/login", strings.NewReader(form.Encode()))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(Login)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if strings.TrimSpace(rr.Body.String()) != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	SetupTestTemplates(t)
	req, err := http.NewRequest(http.MethodGet, "/api/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Logout)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"statusCode":200,"message":"Logout successful"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
