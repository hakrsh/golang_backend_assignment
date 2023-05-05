package msgqueue

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang_backend_assignment/consumer/database"
	"github.com/golang_backend_assignment/consumer/imageutils"
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
	fmt.Println("Successfully Connected to RabbitMQ Instance")
	return conn, nil
}

// NewChannel creates a new amqp.Channel object and returns a pointer to it
func NewChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}
	fmt.Println("Successfully Created a Channel")
	return ch, nil
}

func Consumer(ch *amqp.Channel, queue string, db *sql.DB, image_quality int) {
	msgs, _ := ch.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		fmt.Println("Listening for messages on queue: ", queue)
		for d := range msgs {
			product_id_str := string(d.Body)
			fmt.Println("Received message: ", product_id_str)
			product_id, err := strconv.Atoi(product_id_str)
			if err != nil {
				log.Fatal(err)
				return
			}
			image_urls, err := database.GetProductImages(product_id, db)
			if err != nil {
				log.Fatal(err)
				return
			}
			err = imageutils.DownloadResizeCompressSaveImages(image_urls, image_quality, product_id_str)
			if err != nil {
				log.Fatal(err)
				return
			}
		}
	}()

	<-forever
}
