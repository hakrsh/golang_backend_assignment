package database

import (
	"os"
	"testing"
)

func TestNewDB(t *testing.T) {
	// Set up environment variables for test database
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASSWORD", "example")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_NAME", "product_catalog_db")

	// Call NewDB function and check if there is an error
	db, err := NewDB()
	if err != nil {
		t.Errorf("NewDB returned an error: %v", err)
	}
	defer db.Close()

	// Check if the connection is alive
	err = db.Ping()
	if err != nil {
		t.Errorf("Could not establish a connection to the database: %v", err)
	}
}
