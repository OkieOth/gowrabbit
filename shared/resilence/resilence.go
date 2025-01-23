package resilence

import (
	"fmt"
	"time"
)

type HardenedFuc func() error

func ResilentCall(fn HardenedFuc, maxTries int, millisToWait int, funcName string) (error, int) {
	var lastError error
	for i := 0; i < int(maxTries); i++ {
		if err := fn(); err != nil {
			fmt.Printf("error while execute '%s': %v\n", funcName, err)
			time.Sleep(time.Millisecond * time.Duration(millisToWait*(i+1)))
			lastError = err
		} else {
			fmt.Printf("successfully called '%s'\n", funcName)
			return nil, i
		}
	}
	return fmt.Errorf("finally failed to execute '%s': %v", funcName, lastError), maxTries
}
