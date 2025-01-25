package observer_test

import (
	"sync"
	"testing"
	"time"

	"github.com/okieoth/gowrabbit/shared/observer"
)

func TestObserver(t *testing.T) {
	t.Run("Register_Success", func(t *testing.T) {
		o := observer.NewObserver[int]()
		ch := make(chan int)

		err := o.Register(ch)
		if err != nil {
			t.Errorf("Expected nil, got error: %v", err)
		}

		if o.CountListener() != 1 {
			t.Errorf("Expected 1 listener, got: %d", o.CountListener())
		}
	})

	t.Run("Register_Duplicate", func(t *testing.T) {
		o := observer.NewObserver[int]()
		ch := make(chan int)

		_ = o.Register(ch)
		err := o.Register(ch)
		if err == nil {
			t.Error("Expected error for duplicate registration, got nil")
		}
	})

	t.Run("Unregister_Success", func(t *testing.T) {
		o := observer.NewObserver[int]()
		ch := make(chan int)

		_ = o.Register(ch)
		err := o.Unregister(ch)
		if err != nil {
			t.Errorf("Expected nil, got error: %v", err)
		}

		if o.CountListener() != 0 {
			t.Errorf("Expected 0 listeners, got: %d", o.CountListener())
		}
	})

	t.Run("Unregister_NotFound", func(t *testing.T) {
		o := observer.NewObserver[int]()
		ch := make(chan int)

		err := o.Unregister(ch)
		if err == nil {
			t.Error("Expected error for unregistering non-existent listener, got nil")
		}
	})

	t.Run("Notify_AllListeners", func(t *testing.T) {
		o := observer.NewObserver[int]()
		ch1 := make(chan int, 1)
		ch2 := make(chan int, 1)

		_ = o.Register(ch1)
		_ = o.Register(ch2)

		o.Notify(42)

		select {
		case msg := <-ch1:
			if msg != 42 {
				t.Errorf("Expected 42, got: %d", msg)
			}
		case <-time.After(time.Second):
			t.Error("Timeout waiting for message on ch1")
		}

		select {
		case msg := <-ch2:
			if msg != 42 {
				t.Errorf("Expected 42, got: %d", msg)
			}
		case <-time.After(time.Second):
			t.Error("Timeout waiting for message on ch2")
		}
	})

	t.Run("Notify_PanicRecovery", func(t *testing.T) {
		o := observer.NewObserver[int]()
		ch1 := make(chan int, 1)
		ch2 := make(chan int) // This channel will cause a panic due to no buffer
		close(ch2)
		_ = o.Register(ch1)
		_ = o.Register(ch2)

		o.Notify(42)

		select {
		case msg := <-ch1:
			if msg != 42 {
				t.Errorf("Expected 42, got: %d", msg)
			}
		case <-time.After(time.Second):
			t.Error("Timeout waiting for message on ch1")
		}

		if o.CountListener() != 1 {
			t.Errorf("Expected 1 listener after panic recovery, got: %d", o.CountListener())
		}
	})

	t.Run("Concurrent_RegisterAndNotify", func(t *testing.T) {
		o := observer.NewObserver[int]()
		wg := sync.WaitGroup{}
		wg.Add(3)

		go func() {
			defer wg.Done()
			ch := make(chan int, 1)
			_ = o.Register(ch)
		}()

		go func() {
			defer wg.Done()
			ch := make(chan int, 1)
			_ = o.Register(ch)
		}()

		go func() {
			defer wg.Done()
			o.Notify(99)
		}()

		wg.Wait()
	})
}
