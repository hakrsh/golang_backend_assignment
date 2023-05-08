package msgqueue

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
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
		logrus.Fatalf("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}
	logrus.Info("Successfully Connected to RabbitMQ Instance")
	return conn, nil
}

// NewChannel creates a new amqp.Channel object and returns a pointer to it
func NewChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatalf("Failed to open a channel: %v", err)
		return nil, err
	}
	logrus.Info("Successfully Created a Channel")
	return ch, nil
}

// Take an integer productID and a string queue name and rmq channel as arguments and publish the productID to the queue
func Producer(productID int64, ch *amqp.Channel, queue string) error {
	_, err := ch.QueueDeclare(
		queue, // queue name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		logrus.Errorf("Failed to declare a queue: %v", err)
		return err
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
		logrus.Errorf("Failed to publish a message: %v", err)
		return err
	}
	logrus.Infof("Successfully published productID: %d to queue: %s", productID, queue)
	return err
}
