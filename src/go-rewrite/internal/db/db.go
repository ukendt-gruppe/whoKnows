package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const DATABASE_PATH = "./whoknows.db"

// ConnectDB returns a new connection to the database.
func ConnectDB() (*sql.DB, error) {
	// Attempt to open a connection to the SQLite database
	db, err := sql.Open("sqlite3", DATABASE_PATH)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return db, nil
}

// CheckDBExists checks if the database file exists.
func CheckDBExists() bool {
	if _, err := os.Stat(DATABASE_PATH); os.IsNotExist(err) {
		log.Printf("Database not found at: %s", DATABASE_PATH)
		return false
	}
	return true
}

// InitDB initializes the database with the given schema.
func InitDB(schemaPath string) error {
	// Connect to the database
	db, err := ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close() // Ensure the database connection is closed when this function exits

	// Read the SQL schema file
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %v", err)
	}

	// Execute the schema file content to initialize the database tables
	_, err = db.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %v", err)
	}

	log.Println("Initialized the database:", DATABASE_PATH)
	return nil
}

// QueryDB executes a query and returns the results as a slice of maps.
func QueryDB(query string, args ...interface{}) ([]map[string]interface{}, error) {
	// Connect to the database
	db, err := ConnectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Execute the query
	rows, err := db.Query(query, args...)
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
	// Connect to the database
	db, err := ConnectDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	// Prepare the query
	var id int
	query := "SELECT id FROM users WHERE username = ?"
	err = db.QueryRow(query, username).Scan(&id)

	if err == sql.ErrNoRows {
		return 0, nil // No matching user found
	} else if err != nil {
		return 0, fmt.Errorf("failed to query user ID: %v", err) // An error occurred during the query
	}

	return id, nil
}