package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/swagger"
	"github.com/golang_backend_assignment/producer/db"
	_ "github.com/golang_backend_assignment/producer/docs"
	"github.com/golang_backend_assignment/producer/handlers"
	"github.com/golang_backend_assignment/producer/msgqueue"
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

	queue := os.Getenv("RM_QUEUENAME")
	conn, err := msgqueue.NewRMQ()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	ch, err := msgqueue.NewChannel(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	// Create the Fiber app
	app := fiber.New()

	// Define the route to receive the product data
	app.Post("/products", handlers.SaveProduct(db, ch, queue))
	app.Get("/swagger/*", swagger.HandlerDefault)
	// Start the server
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
