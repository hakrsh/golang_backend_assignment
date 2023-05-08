package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"reflect"
	"testing"
)

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
	// Call the GetProductImages function with product_id = 999
	images, err = GetProductImages(999, db)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}

}

func TestUpdateProductImages(t *testing.T) {
	// Set up a test database
	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening test database: %v", err)
	}
	defer testDB.Close()

	// Create the Products table
	_, err = testDB.Exec(`
		CREATE TABLE Products (
			product_id INTEGER PRIMARY KEY AUTOINCREMENT,
			product_name TEXT,
			product_description TEXT,
			product_images TEXT,
			compressed_product_images TEXT,
			product_price REAL,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("Error creating Products table: %v", err)
	}

	// Insert a test product
	_, err = testDB.Exec(`
		INSERT INTO Products (product_name, product_description, product_images, product_price, created_at)
		VALUES ('Test Product', 'This is a test product.', 'image1.jpg', 19.99, datetime('now'))
	`)
	if err != nil {
		t.Fatalf("Error inserting test product: %v", err)
	}

	// Update the product images
	productID := 1
	err = UpdateProductImages(testDB, productID, []string{"image1.jpg", "image2.jpg"})
	if err != nil {
		t.Fatalf("Error updating product images: %v", err)
	}

	// Check that the images were updated correctly
	var compressedImages string
	err = testDB.QueryRow("SELECT compressed_product_images FROM Products WHERE product_id = ?", productID).Scan(&compressedImages)
	if err != nil {
		t.Fatalf("Error getting compressed product images: %v", err)
	}
	expectedImages := "image1.jpg,image2.jpg"
	if compressedImages != expectedImages {
		t.Fatalf("Expected compressed images to be %q, but got %q", expectedImages, compressedImages)
	}

	// Update the product images again, this time with an empty list
	err = UpdateProductImages(testDB, productID, []string{})
	if err != nil {
		t.Fatalf("Error updating product images: %v", err)
	}

	// Check that the images were not updated
	err = testDB.QueryRow("SELECT compressed_product_images FROM Products WHERE product_id = ?", productID).Scan(&compressedImages)
	if err != nil {
		t.Fatalf("Error getting compressed product images: %v", err)
	}
	if compressedImages != expectedImages {
		t.Fatalf("Expected compressed images to be %q, but got %q", expectedImages, compressedImages)
	}

	// Update a nonexistent product ID
	err = UpdateProductImages(testDB, 999, []string{"image3.jpg"})
	if err == nil {
		t.Fatal("Expected error updating nonexistent product, but got nil")
	}
}
