package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"reflect"
	"testing"
)

func TestNewDB(t *testing.T) {
	// Set up environment variables for test database
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASSWORD", "example")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_NAME", "ProductsAndUsersDB")

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

func TestGetProductImages(t *testing.T) {
	// Open a test database connection
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to open database: %s", err)
	}
	defer db.Close()

	// Create a test Products table with some data
	_, err = db.Exec(`
        CREATE TABLE Products (
            product_id INTEGER PRIMARY KEY,
            product_images TEXT
        );
        INSERT INTO Products (product_id, product_images)
        VALUES (1, "image1.jpg,image2.jpg,image3.jpg");
    `)
	if err != nil {
		t.Fatalf("Failed to create test table: %s", err)
	}

	// Call the GetProductImages function with product_id = 1
	images, err := GetProductImages(1, db)
	if err != nil {
		t.Fatalf("Failed to get product images: %s", err)
	}

	// Check that the result is as expected
	expected := []string{"image1.jpg", "image2.jpg", "image3.jpg"}
	if !reflect.DeepEqual(images, expected) {
		t.Errorf("Unexpected result: got %v, expected %v", images, expected)
	}
}
