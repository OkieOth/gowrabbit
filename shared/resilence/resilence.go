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
			fmt.Printf("[%d/%d] error while execute '%s': %v\n", i+1, maxTries, funcName, err)
			if i < maxTries {
				time.Sleep(time.Millisecond * time.Duration(millisToWait*(i+1)))
				fmt.Printf("[%d/%d] ... awake again ...\n", i+1, maxTries)
			}
			lastError = err
		} else {
			fmt.Printf("[%d/%d] successfully called '%s'\n", i+1, maxTries, funcName)
			return nil, i + 1
		}
	}
	fmt.Printf("[%d/%d] ... done :-/\n", maxTries, maxTries)
	return fmt.Errorf("[%d/%d] finally failed to execute '%s': %v", maxTries, maxTries, funcName, lastError), maxTries
}
