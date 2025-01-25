package amqp_test

import (
	"testing"

	"github.com/okieoth/gowrabbit/shared/amqp"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	s1 := amqp.NewServer()
	assert.Equal(t, "localhost", s1.Host)
	assert.Equal(t, uint(5672), s1.Port)

	s2 := amqp.NewServer(
		amqp.Host("test.com"),
	)
	assert.Equal(t, "test.com", s2.Host)
	assert.Equal(t, uint(5672), s1.Port)

	s3 := amqp.NewServer(
		amqp.Port(8000),
	)
	assert.Equal(t, "localhost", s3.Host)
	assert.Equal(t, uint(8000), s3.Port)

	s4 := amqp.NewServer(
		amqp.Host("test2.com"),
		amqp.Port(8001),
	)
	assert.Equal(t, "test2.com", s4.Host)
	assert.Equal(t, uint(8001), s4.Port)
}
