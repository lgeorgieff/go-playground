package main

import (
	"fmt"
	"time"
)

func NewOrDoneChan[T any](stream <-chan T, done <-chan struct{}) <-chan T {
	orStream := make(chan T)

	go func() {
		defer close(orStream)

		for {
			select {
			case <-done:
				return
			case t, ok := <-stream:
				if !ok {
					return
				}
				select {
				case <-done:
				case orStream <- t:
				}
			}
		}
	}()

	return orStream
}

func main() {
	const (
		useDoneChan     = false
		closeStreamChan = false
	)

	stream := make(chan int)
	done := make(chan struct{})

	go func() {
		defer func() {
			if !closeStreamChan {
				close(stream)
			}
		}()

		for i := 0; i < 10; i++ {
			time.Sleep(100 * time.Millisecond)

			if closeStreamChan && i == 6 {
				close(stream)
				return
			}

			if useDoneChan && i == 7 {
				close(done)
			}

			stream <- i
		}
	}()

	for i := range NewOrDoneChan(stream, done) {
		fmt.Println(i)
	}
}
