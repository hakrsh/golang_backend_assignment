package msgqueue

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

// NewRMQ sets up a connection to a RabbitMQ server and returns a pointer to an amqp.Connection object
func NewRMQ() (*amqp.Connection, error) {
	rmqHost := os.Getenv("RMQ_HOST")
	rmqPort := os.Getenv("RMQ_PORT")
	rmqUser := os.Getenv("RMQ_USER")
	rmqPassword := os.Getenv("RMQ_PASSWORD")

	rmqURL := fmt.Sprintf("amqp://%s:%s@%s:%s/", rmqUser, rmqPassword, rmqHost, rmqPort)
	conn, err := amqp.Dial(rmqURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	fmt.Println("Successfully Connected to our RabbitMQ Instance")
	return conn, nil
}
