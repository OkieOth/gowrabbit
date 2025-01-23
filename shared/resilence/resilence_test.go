package resilence_test

import (
	"fmt"
	"testing"

	"github.com/okieoth/gowrabbit/shared/resilence"
)

func TestFail(t *testing.T) {
	callCount := 0
	succedAtFunc := func(succedAt int) bool {
		callCount++
		if callCount == succedAt {
			return true
		} else {
			return false
		}
	}

	testFn := func() error {
		if !succedAtFunc(4) {
			return fmt.Errorf("Failed in call: %d", callCount)
		} else {
			return nil
		}
	}

	if err, tryCount := resilence.ResilentCall(testFn, 3, 100, "test-fail"); err == nil {
		t.Errorf("function didn't fail 'test-fail', tries: %d", tryCount)
	} else {
		if tryCount != 3 {
			t.Errorf("function 'test-fail' failed after wrong tries, received tries: %d", tryCount)
		}
	}
}

func TestSucceed(t *testing.T) {
	callCount := 0
	succedAtFunc := func(succedAt int) bool {
		callCount++
		if callCount == succedAt {
			return true
		} else {
			return false
		}
	}

	testFn := func() error {
		if !succedAtFunc(4) {
			return fmt.Errorf("Failed in call: %d", callCount)
		} else {
			return nil
		}
	}

	if err, tryCount := resilence.ResilentCall(testFn, 4, 100, "test-succeed"); err == nil {
		if tryCount != 3 {
			t.Errorf("function 'test-succeed' successful executed, but after wrong tries, received tries: %d", tryCount)
		}
	} else {
		t.Errorf("function 'test-succeed' failed, tries: %d", tryCount)
	}
}
