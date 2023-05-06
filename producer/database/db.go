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

func UserExists(db *sql.DB, userID int) error {
	userStmt, err := db.Prepare("SELECT COUNT(*) FROM Users WHERE id = ?")
	if err != nil {
		return err
	}
	defer userStmt.Close()

	var count int
	err = userStmt.QueryRow(userID).Scan(&count)
	if err != nil {
		return err
	}

	return nil
}

func InsertProduct(db *sql.DB, ProductName string, ProductDescription string, ProductPrice float64, productImages []string) (int64, error) {
	// Join the product images into a comma-separated string
	productImagesStr := strings.Join(productImages, ",")

	// Insert the product into the database
	stmt, err := db.Prepare("INSERT INTO Products (product_name, product_description, product_images, product_price, created_at) VALUES (?, ?, ?, ?, NOW())")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(ProductName, ProductDescription, productImagesStr, ProductPrice)
	if err != nil {
		return 0, err
	}

	// Get the product ID
	productID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return productID, nil
}
