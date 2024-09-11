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