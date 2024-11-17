// File: src/backend/internal/handlers/api_test.go

package handlers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/sessions"
	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/db"
	"golang.org/x/crypto/bcrypt"
)

var testStore = sessions.NewCookieStore([]byte("test-secret-key"))

func setupTest(t *testing.T) (sqlmock.Sqlmock, func()) {
	// Create SQL mock
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}

	// Set up our mocked database
	db.DB = mockDB

	// Mock templates
	templates = template.Must(template.New("mock").Parse("mock template"))

	// Return cleanup function
	return mock, func() {
		mockDB.Close()
	}
}

func setupTestSession(r *http.Request) {
	session, _ := testStore.New(r, "session-name")
	ctx := context.WithValue(r.Context(), "session", session)
	*r = *r.WithContext(ctx)
}

func TestRegister(t *testing.T) {
	mock, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name           string
		method         string
		username       string
		email          string
		password       string
		expectedStatus int
		expectedBody   string
		mockSetup      func(sqlmock.Sqlmock)
	}{
		{
			name:           "Valid registration",
			method:         http.MethodPost,
			username:       "testuser",
			email:          "test@example.com",
			password:       "testpass",
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"statusCode":201,"message":"User registered successfully"}`,
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Expect check for existing user
				mock.ExpectQuery(`SELECT (.+) FROM users WHERE username = \$1`).
					WithArgs("testuser").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password"}))

				// Expect insert new user
				mock.ExpectExec(`INSERT INTO users`).
					WithArgs("testuser", "test@example.com", sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:           "Username already exists",
			method:         http.MethodPost,
			username:       "existinguser",
			email:          "test@example.com",
			password:       "testpass",
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"statusCode":409,"message":"Username already exists"}`,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM users WHERE username = \$1`).
					WithArgs("existinguser").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password"}).
						AddRow(1, "existinguser", "test@example.com", "hashedpass"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock expectations
			tt.mockSetup(mock)

			// Create request
			form := url.Values{}
			form.Add("username", tt.username)
			form.Add("email", tt.email)
			form.Add("password", tt.password)

			req, err := http.NewRequest(tt.method, "/api/register", strings.NewReader(form.Encode()))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			setupTestSession(req)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(Register)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if strings.TrimSpace(rr.Body.String()) != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt.expectedBody)
			}

			// Verify all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	mock, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name           string
		query          string
		language       string
		expectedStatus int
		mockSetup      func(sqlmock.Sqlmock)
	}{
		{
			name:           "Valid search",
			query:          "test",
			language:       "en",
			expectedStatus: http.StatusOK,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"title", "url", "content", "source"}).
					AddRow("Test Title", "http://test.com", "Test Content", "page")
				mock.ExpectQuery(`SELECT (.+) FROM pages`).
					WithArgs("en", "%test%").
					WillReturnRows(rows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			url := fmt.Sprintf("/api/search?q=%s&language=%s", tt.query, tt.language)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(Search)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	// Add mock database setup
	mock, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name   string
		method string

		username       string
		password       string
		expectedStatus int
		expectedBody   string
		mockSetup      func(sqlmock.Sqlmock)
	}{
		{
			name:           "Valid login",
			method:         http.MethodPost,
			username:       "testuser",
			password:       "testpass",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"statusCode":200,"message":"Login successful"}`,
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Create an actual hashed password for "testpass"
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password"}).
					AddRow(1, "testuser", "test@example.com", string(hashedPassword))
				mock.ExpectQuery(`SELECT (.+) FROM users WHERE username = \$1`).
					WithArgs("testuser").
					WillReturnRows(rows)
			},
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
			// Set up mock expectations if they exist
			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			form := url.Values{}
			form.Add("username", tt.username)
			form.Add("password", tt.password)

			req, err := http.NewRequest(tt.method, "/api/login", strings.NewReader(form.Encode()))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			setupTestSession(req)

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

	setupTestSession(req)

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
