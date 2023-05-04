package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type Product struct {
	UserID             int      `json:"user_id"`
	ProductName        string   `json:"product_name"`
	ProductDescription string   `json:"product_description"`
	ProductImages      []string `json:"product_images"`
	ProductPrice       float64  `json:"product_price"`
}

// @Summary Save a product
// @Description Save a product to the database
// @Tags Products
// @Accept json
// @Produce json
// @Param product body Product true "Product data"
// @Success 200 {string} string "Product saved successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Internal server error"
// @Router /products [post]
func SaveProduct(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse the request body into a Product struct
		var product Product
		if err := c.BodyParser(&product); err != nil {
			return err
		}
		// Convert the product images slice to a comma-separated string
		productImagesStr := ""
		for _, image := range product.ProductImages {
			productImagesStr += image + ","
		}
		productImagesStr = productImagesStr[:len(productImagesStr)-1]
		// Insert the product into the database
		stmt, err := db.Prepare("INSERT INTO Products (product_name, product_description, product_images, product_price, created_at) VALUES (?, ?, ?, ?, NOW())")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(product.ProductName, product.ProductDescription, productImagesStr, product.ProductPrice)
		if err != nil {
			return err
		}

		// Return a success message
		return c.SendString("Product saved successfully")
	}
}
