package amqp_test

import (
	"testing"

	"github.com/okieoth/gowrabbit/shared/amqp"
)

func TestNewServer(t *testing.T) {
	s1 := amqp.NewServer()

	s2 := amqp.NewServer(
		amqp.Host("test.com"),
	)

	s3 := amqp.NewServer(
		amqp.Port(8000),
	)

	s4 := amqp.NewServer(
		amqp.Host("test.com"),
		amqp.Port(8000),
	)
}

func TestAddConnNotify(t *testing.T) {
	// Initialize a new Connection instance
	conn := amqp.NewConnection()

	// Create a channel to add
	notifyChan := make(chan amqp.ConnectionState)

	// Call AddConnNotify
	conn.AddConnNotify(notifyChan)

	// Verify that the channel is added
	if conn.ConnNotifyCount() != 1 {
		t.Errorf("AddConnNotify failed: expected 1, got %v", conn.ConnNotifyCount())
	}

	go func() {
		conn.SendConnNotify(amqp.DISCONNECTED)
	}()
	// TODO reading with timeout
}

func TestDelConnNotify(t *testing.T) {
	// Initialize a new Connection instance
	conn := amqp.NewConnection()

	// Create channels to add and then remove
	notifyChan1 := make(chan amqp.ConnectionState)
	notifyChan2 := make(chan amqp.ConnectionState)

	// Add channels
	conn.AddConnNotify(notifyChan1)
	conn.AddConnNotify(notifyChan2)

	if conn.ConnNotifyCount() != 2 {
		t.Fatalf("Setup failed: expected 2 channels, got %d", conn.ConnNotifyCount())
	}

	// Remove one channel
	conn.DelConnNotify(notifyChan1)

	// Verify that the channel is removed
	if conn.ConnNotifyCount() != 1 {
		t.Errorf("DelConnNotify failed: expected 1, got %v", conn.ConnNotifyCount())
	}

	// Remove the remaining channel
	conn.DelConnNotify(notifyChan2)

	// Verify that all channels are removed
	if conn.ConnNotifyCount() != 0 {
		t.Errorf("DelConnNotify failed: expected no channels, got %d", conn.ConnNotifyCount())
	}
}
