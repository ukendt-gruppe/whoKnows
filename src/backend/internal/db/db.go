package db

import (
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/utils"
)

var DB *sql.DB

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidUser  = errors.New("invalid user")
)

func ConnectDB() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	var err error
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Printf("Successfully connected to PostgreSQL database")
	return nil
}

// QueryDB executes a query and returns the results as a slice of maps.
func QueryDB(query string, args ...interface{}) ([]map[string]interface{}, error) {
	// Execute the query
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	// Get the column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get column names: %v", err)
	}

	// Prepare a slice of maps to hold the query results
	var results []map[string]interface{}
	for rows.Next() {
		// Create a slice of interfaces to hold each column value
		columnsData := make([]interface{}, len(columns))
		columnsPointers := make([]interface{}, len(columns))
		for i := range columnsData {
			columnsPointers[i] = &columnsData[i]
		}

		// Scan the row data into column pointers
		if err := rows.Scan(columnsPointers...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		// Create a map to store the row data
		rowMap := make(map[string]interface{})
		for i, colName := range columns {
			rowMap[colName] = columnsData[i]
		}

		// Add the row map to the results
		results = append(results, rowMap)
	}

	return results, nil
}

// GetUserID looks up the ID for a given username.
func GetUserID(username string) (int, error) {
	var id int
	// Note: Changed ? to $1 for PostgreSQL parameterization
	query := "SELECT id FROM users WHERE username = $1"
	err := DB.QueryRow(query, username).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil // No matching user found
	} else if err != nil {
		return 0, fmt.Errorf("failed to query user ID: %v", err)
	}
	return id, nil
}

// GetUser retrieves a user by username or ID.
func GetUser(identifier interface{}) (*User, error) {
	var user User
	var err error
	switch v := identifier.(type) {
	case string:
		// Note: Changed ? to $1 for PostgreSQL parameterization
		err = DB.QueryRow("SELECT id, username, email, password, needs_password_reset FROM users WHERE username = $1", v).Scan(
			&user.ID, &user.Username, &user.Email, &user.Password, &user.NeedsPasswordReset)
	case int:
		err = DB.QueryRow("SELECT id, username, email, password, needs_password_reset FROM users WHERE id = $1", v).Scan(
			&user.ID, &user.Username, &user.Email, &user.Password, &user.NeedsPasswordReset)
	default:
		return nil, ErrInvalidUser
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user in the database.
func CreateUser(username, email, password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	// Note: Changed ? to $1, $2, $3 for PostgreSQL parameterization
	_, err = DB.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)",
		username, email, hashedPassword)
	return err
}

// User represents a user in the database.
type User struct {
	ID                int
	Username          string
	Email             string
	Password          string
	NeedsPasswordReset bool
}

func init() {
	gob.Register(&User{})
}
