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
