package main

import (
	"fmt"

	"time"

	"github.com/okieoth/gowrabbit/pub"
	"github.com/okieoth/gowrabbit/shared/amqp"
	"github.com/okieoth/gowrabbit/sub"
)

func main() {
	fmt.Println("Hello from main :)")
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
		fmt.Println("Successfully connected :)")
		for {
			fmt.Println("I am going to sleep for 10s ...")
			time.Sleep(10 * time.Second)
		}
	} else {
		fmt.Println("Connection failed :-/ ")
		fmt.Println(err)
	}
}
