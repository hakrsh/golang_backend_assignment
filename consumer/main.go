package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang_backend_assignment/consumer/db"
	"github.com/golang_backend_assignment/consumer/msgqueue"
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
	msgqueue.Consumer(ch, queue)
}
