package main

import (
	"fmt"

	"github.com/okieoth/gowrabbit/pub"
	"github.com/okieoth/gowrabbit/sub"
)

func main() {
	fmt.Println("Hello from main :)")
	pub.DummyPub()
	sub.DummySub()
}
