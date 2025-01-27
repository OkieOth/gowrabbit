package tests

import (
	"testing"

	"github.com/okieoth/gowrabbit/shared/amqp"
	"github.com/stretchr/testify/assert"
)

func TestConnectionFailure_IT(t *testing.T) {
	connection := amqp.NewConnection(
		amqp.User("guest"),
		amqp.Password("guest"),
		amqp.Servers([]amqp.Server{
			amqp.NewServer(
				amqp.Host("not.existing.host"),
			),
		}),
		amqp.ResilenceMaxRetries(5),
		amqp.ResilenceWaitMilis(100),
	)
	//fmt.Printf("Connection: %v\n", connection)
	err := connection.Connect()
	assert.NotNil(t, err)
}

func TestConnection_IT(t *testing.T) {
	user, password, host, port := getConnParamsToUse()

	connection := amqp.NewConnection(
		amqp.User(user),
		amqp.Password(password),
		amqp.Servers([]amqp.Server{
			amqp.NewServer(
				amqp.Host(host),
				amqp.Port(uint(port)),
			),
		}),
	)
	//fmt.Printf("Connection: %v\n", connection)
	err := connection.Connect()
	assert.Nil(t, err)
}
