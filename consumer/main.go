package main

import (
	"os"

	"github.com/golang_backend_assignment/consumer/database"
	"github.com/golang_backend_assignment/consumer/msgqueue"
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
	image_quality := 60
	msgqueue.Consumer(ch, queue, db, image_quality)
}
