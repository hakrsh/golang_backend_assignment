package main

import (
	"os"

	fiber "github.com/gofiber/fiber/v2"

	"github.com/gofiber/swagger"
	"github.com/golang_backend_assignment/producer/database"
	_ "github.com/golang_backend_assignment/producer/docs"
	"github.com/golang_backend_assignment/producer/handlers"
	"github.com/golang_backend_assignment/producer/msgqueue"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Error("Error loading .env file")
		return
	}
	// Connect to the database
	db, err := database.NewDB()
	if err != nil {
		logrus.Errorf("Failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	queue := os.Getenv("RM_QUEUENAME")
	conn, err := msgqueue.NewRMQ()
	if err != nil {
		logrus.Errorf("Failed to connect to RabbitMQ: %v", err)
		return
	}
	defer conn.Close()
	ch, err := msgqueue.NewChannel(conn)
	if err != nil {
		logrus.Errorf("Failed to open a rmq channel: %v", err)
		return
	}
	defer ch.Close()

	// Create the Fiber app
	app := fiber.New()

	// Define the route to receive the product data
	app.Post("/products", handlers.SaveProduct(db, ch, queue))
	app.Get("/swagger/*", swagger.HandlerDefault)
	// Start the server
	if err := app.Listen(":3000"); err != nil {
		logrus.Fatalf("Error in starting the server...: %v", err)
	}
}
