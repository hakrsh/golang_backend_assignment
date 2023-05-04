package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/swagger"
	"github.com/golang_backend_assignment/db"
	_ "github.com/golang_backend_assignment/docs"
	"github.com/golang_backend_assignment/handlers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	// Connect to the database
	db, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the Fiber app
	app := fiber.New()

	// Define the route to receive the product data
	app.Post("/products", handlers.SaveProduct(db))
	app.Get("/swagger/*", swagger.HandlerDefault)
	// Start the server
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
