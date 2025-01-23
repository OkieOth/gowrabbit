package amqp

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Chann struct {
	connection *Connection
	channel    *amqp.Channel
}

func NewChannel(conn *Connection) (Chann, error) {
	ret := Chann{
		connection: conn,
	}

	c, err := ret.connection.conn.Channel()
	if err != nil {
		return ret, fmt.Errorf("error while creating a channel: %v", err)
	}
	ret.channel = c
	go func() {
		channelClosedChan := make(chan *amqp.Error)
		if ret.channel != nil {
			err := ret.channel.NotifyClose(channelClosedChan)
			if err != nil {
				// TODO logging
			} else {
				if e, ok := <-channelClosedChan; ok {
					fmt.Println("channel was closed w/ error: ", e)
				} else {
					fmt.Println("channel closed")
				}
			}

		} else {
			// TODO logging
			fmt.Println("connection object is nil")
		}
	}()

	return ret, nil
}
