package tests

import (
	"fmt"
	"testing"
)

// In case you have installed
func TestIsRabbitmqAdminLocalAvailable(t *testing.T) {
	b := isRabbitmqAdminLocalAvailable()
	fmt.Println("rabbitmqadmin is locally available: ", b)
}
