package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/golang_backend_assignment/producer/database"
	"github.com/golang_backend_assignment/producer/msgqueue"
	"github.com/streadway/amqp"
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
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /products [post]
func SaveProduct(db *sql.DB, ch *amqp.Channel, queue string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse the request body into a Product struct
		var product Product
		if err := c.BodyParser(&product); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
		}

		err := database.UserExists(db, product.UserID)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}

		productID, err := database.InsertProduct(db, product.ProductName, product.ProductDescription, product.ProductPrice, product.ProductImages)
		err = msgqueue.Producer(productID, ch, queue)
		if err != nil {
			return err
		}
		// Return a success message
		return c.SendString("Product saved successfully")
	}
}
