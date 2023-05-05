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

// Take an integer productID and a string queue name and rmq channel as arguments and publish the productID to the queue
func Producer(productID int64, ch *amqp.Channel, queue string) error {
	_, err := ch.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		fmt.Println(err)
	}

	err = ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("%d", productID)),
		},
	)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Published Message to Queue")
	return err
}

func Consumer(ch *amqp.Channel, queue string) {
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
			fmt.Println("Received message: ", string(d.Body))
		}
	}()

	<-forever
}
