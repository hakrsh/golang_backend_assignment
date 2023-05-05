package msgqueue

import (
	"os"
	"testing"
)

func TestNewRMQ(t *testing.T) {
	os.Setenv("RMQ_HOST", "localhost")
	os.Setenv("RMQ_PORT", "5672")
	os.Setenv("RMQ_USER", "guest")
	os.Setenv("RMQ_PASSWORD", "guest")

	conn, err := NewRMQ()
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Verify that the connection is not nil
	if conn == nil {
		t.Error("Expected connection to not be nil")
	}
}
