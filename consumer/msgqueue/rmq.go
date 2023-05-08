package msgqueue

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/golang_backend_assignment/consumer/database"
	"github.com/golang_backend_assignment/consumer/imageutils"
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
		logrus.Errorf("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}
	logrus.Info("Successfully Connected to RabbitMQ Instance")
	return conn, nil
}

// NewChannel creates a new amqp.Channel object and returns a pointer to it
func NewChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		logrus.Errorf("Failed to open a channel: %v", err)
		return nil, err
	}
	logrus.Info("Successfully Created a Channel")
	return ch, nil
}

func Consumer(ch *amqp.Channel, queue string, db *sql.DB, image_quality int) {
	_, err := ch.QueueDeclare(
		queue, // queue name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		logrus.Errorf("Failed to declare queue: %v", err)
		return
	}
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
		logrus.Info("Listening for messages on queue: ", queue)
		for d := range msgs {
			product_id_str := string(d.Body)
			logrus.Info("Received message: ", product_id_str)
			// Spawn a new goroutine to process the message
			go func(product_id_str string) {
				product_id, err := strconv.Atoi(product_id_str)
				if err != nil {
					logrus.Errorf("Failed to convert product_id to int: %v", err)
					return
				}
				image_urls, err := database.GetProductImages(product_id, db)
				if err != nil {
					logrus.Errorf("Error in fetching product images from db: %v", err)
					return
				}
				err, compressedImagePaths := imageutils.DownloadResizeCompressSaveImages(image_urls, image_quality, product_id_str)
				if err != nil {
					logrus.Errorf("Error in DownloadResizeCompressSaveImages: %v", err)
					return
				}
				err = database.UpdateProductImages(db, product_id, compressedImagePaths)
				if err != nil {
					logrus.Errorf("Error in updating product images in db: %v", err)
					return
				}
			}(product_id_str)
		}
	}()

	<-forever
}
