package observer

import (
	"fmt"
	"sync"
)

type Observer[T any] struct {
	listener []chan<- T
	mutex    sync.Mutex
}

func NewObserver[T any]() Observer[T] {
	return Observer[T]{
		listener: make([]chan<- T, 0),
	}
}

func (o *Observer[T]) Register(listener chan<- T) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	for i := 0; i < len(o.listener); i++ {
		if o.listener[i] == listener {
			return fmt.Errorf("listener already registered")
		}
	}
	o.listener = append(o.listener, listener)
	return nil
}

func (o *Observer[T]) Unregister(listener chan<- T) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	for i := 0; i < len(o.listener); i++ {
		if o.listener[i] == listener {
			o.listener = append(o.listener[:i], o.listener[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("no related listener found")
}

func (o *Observer[T]) CountListener() int {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	return len(o.listener)
}

func (o *Observer[T]) Notify(msg T) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	failedChanIndexes := make([]int, 0)
	for i, ch := range o.listener {
		func(index int, ch chan<- T) {
			defer func() {
				if r := recover(); r != nil {
					failedChanIndexes = append(failedChanIndexes, index)
				}
			}()
			ch <- msg
		}(i, ch)
	}
	for i := len(failedChanIndexes) - 1; i >= 0; i-- {
		index := failedChanIndexes[i]
		o.listener = append(o.listener[:index], o.listener[index+1:]...)
	}
}
