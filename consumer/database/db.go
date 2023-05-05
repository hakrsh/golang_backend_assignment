package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB() (*sql.DB, error) {
	// Load database details from environment file
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Create database connection string
	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Connect to the database
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to the database: ", err)
		return nil, err
	}
	fmt.Println("Successfully connected to the database")
	return db, nil
}

func GetProductImages(product_id int, db *sql.DB) ([]string, error) {
	fmt.Print("Getting product images for product_id: ", product_id)
	stmt, err := db.Prepare("SELECT product_images FROM Products WHERE product_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the SELECT statement
	var product_images string
	err = stmt.QueryRow(product_id).Scan(&product_images)
	if err != nil {
		return nil, err
	}

	// Split the comma-separated values and return them as a slice of strings
	images := strings.Split(product_images, ",")
	return images, nil
}
