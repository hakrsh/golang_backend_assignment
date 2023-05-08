#!/bin/bash
set -e

# Start containers using docker-compose
echo "Starting containers using docker-compose..."
docker-compose up -d

# Wait for containers to start
echo "Waiting for containers to start..."
sleep 10

# Run producer and consumer apps using Go
echo "Starting producer and consumer apps using Go..."
cd ./consumer
go run main.go &
producer_pid=$!
echo "Producer PID: $producer_pid"
cd ../producer
go run main.go &
consumer_pid=$!
echo "Consumer PID: $consumer_pid"

# Wait for producer and consumer to start
echo "Waiting for producer and consumer to start..."
sleep 5

# Send a POST request to the product API
echo "Sending a POST request to the product API..."
curl -X 'POST' \
    'http://localhost:3000/products' \
    -H 'accept: application/json' \
    -H 'Content-Type: application/json' \
    -d '{
      "user_id": 1,
      "product_name": "Headphones",
      "product_description": "This Headphones will blow your mind!",
      "product_images": [
          "https://raw.githubusercontent.com/harikrishnanum/products/main/samsung.jpg",
          "https://raw.githubusercontent.com/harikrishnanum/products/main/jbl.jpg",
          "https://raw.githubusercontent.com/harikrishnanum/products/main/sony.jpeg"
      ],
      "product_price": 10000
  }'

# wait for the consumer to consume all the messages
echo "Waiting for the consumer to consume all the messages..."
sleep 10
# Stop producer and consumer and remove volumes using docker-compose
echo "Stopping producer and consumer apps using Go..."
kill $producer_pid
kill $consumer_pid
echo "Stopped producer and consumer apps using Go."

# Stop the process listening to port 3000
echo "Stopping the process listening to port 3000..."
lsof -ti tcp:3000 | xargs kill
echo "Stopped the process listening to port 3000."

# Stop containers and remove volumes using docker-compose
echo "Stopping containers using docker-compose..."
docker-compose down -v
echo "Stopped containers using docker-compose."
echo "Integration test completed successfully."
