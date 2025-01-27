package main

import (
	"time"

	"github.com/okieoth/gowrabbit/pub"
	"github.com/okieoth/gowrabbit/shared/amqp"
	"github.com/okieoth/gowrabbit/sub"

	"basic/logger"
)

func main() {

	log := logger.NewLogger()
	pub.DummyPub()
	sub.DummySub()

	connection := amqp.NewConnection(
		amqp.User("guest"),
		amqp.Password("guest"),
		amqp.Servers([]amqp.Server{
			amqp.NewServer(
				amqp.Host("localhost"),
			),
		}),
	)
	//fmt.Printf("Connection: %v\n", connection)
	if err := connection.Connect(); err == nil {
		log.Info("Successfully connected :)")
		for {
			log.Info("I am going to sleep for 10s ...")
			time.Sleep(10 * time.Second)
		}
	} else {
		log.Error("Connection failed: %v", err)
	}
}
