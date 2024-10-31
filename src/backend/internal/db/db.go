package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"encoding/gob"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ukendt-gruppe/whoKnows/src/backend/internal/utils"
)

const DATABASE_PATH = "./internal/db/whoknows.db"

var DB *sql.DB

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidUser  = errors.New("invalid user")
)

// InitDB initializes the database connection and creates tables if they don't exist.
func InitDB(schemaPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", DATABASE_PATH)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Check if the database file exists
	if _, err := os.Stat(DATABASE_PATH); os.IsNotExist(err) {
		log.Printf("Database not found at: %s. Creating new database.", DATABASE_PATH)
	}

	// Read and execute the schema
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %v", err)
	}

	_, err = DB.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %v", err)
	}

	log.Println("Initialized the database:", DATABASE_PATH)
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
	query := "SELECT id FROM users WHERE username = ?"
	err := DB.QueryRow(query, username).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil // No matching user found
	} else if err != nil {
		return 0, fmt.Errorf("failed to query user ID: %v", err)
	}
	return id, nil
}

// GetUser retrieves a user by username.
func GetUser(identifier interface{}) (*User, error) {
	var user User
	var err error
	switch v := identifier.(type) {
	case string:
		err = DB.QueryRow("SELECT id, username, email, password FROM users WHERE username = ?", v).Scan(
			&user.ID, &user.Username, &user.Email, &user.Password)
	case int:
		err = DB.QueryRow("SELECT id, username, email, password FROM users WHERE id = ?", v).Scan(
			&user.ID, &user.Username, &user.Email, &user.Password)
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

	_, err = DB.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		username, email, hashedPassword)
	return err
}

// User represents a user in the database.
type User struct {
  ID       int
  Username string
  Email    string
  Password string
}

func init() {
  gob.Register(&User{})
}