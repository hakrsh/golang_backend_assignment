package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
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
		logrus.Errorf("Error connecting to the database: %v", err)
		return nil, err
	}
	logrus.Info("Successfully connected to the database")
	return db, nil
}

func ProductExists(db *sql.DB, productID int) error {
	logrus.Info("Checking if product exists for product_id: ", productID)
	productStmt, err := db.Prepare("SELECT COUNT(*) FROM Products WHERE product_id = ?")
	if err != nil {
		logrus.Errorf("Error preparing SQL statement: %v", err)
		return err
	}
	defer productStmt.Close()

	var count int
	err = productStmt.QueryRow(productID).Scan(&count)
	if err != nil {
		logrus.Errorf("Error executing SQL statement: %v", err)
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func GetProductImages(product_id int, db *sql.DB) ([]string, error) {
	err := ProductExists(db, product_id)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("Product does not exist for product_id: %d", product_id)
			return nil, err
		}
		logrus.Errorf("Error checking if product exists for product_id: %d", product_id)
		return nil, err
	}

	logrus.Info("Getting product images for product_id: ", product_id)
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

func UpdateProductImages(db *sql.DB, productID int, compressedImagesPaths []string) error {
	err := ProductExists(db, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("Product does not exist for product_id: %d", productID)
			return err
		}
		logrus.Errorf("Error checking if product exists for product_id: %d", productID)
		return err
	}
	// Update the database
	if len(compressedImagesPaths) == 0 {
		logrus.Error("No images to update")
		return nil
	}
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	compressedImages := strings.Join(compressedImagesPaths, ",")
	query := "UPDATE Products SET compressed_product_images = ?, updated_at = ? WHERE product_id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		logrus.Errorf("error preparing update statement: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(compressedImages, currentTime, productID)
	if err != nil {
		logrus.Errorf("error executing update statement: %v", err)
		return err
	}
	logrus.Infof("Successfully updated product_id: %d", productID)
	return nil
}
