package utils

import (
	"time"
)

func GetValueWithTimeout[T any](ch chan T, timeout time.Duration) *T {
	timeoutChan := make(chan bool, 1)
	go timeoutFunc(timeoutChan, timeout)
	select {
	case val := <-ch:
		return &val
	case <-timeoutChan:
		return nil
	}
}

func timeoutFunc(timeoutChan chan bool, timeout time.Duration) {
	time.Sleep(timeout)
	timeoutChan <- true
}
