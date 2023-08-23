package main

import (
	"fmt"
	"time"
)

func OrChan(channels ...<-chan struct{}) <-chan struct{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan struct{})
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-OrChan(append(channels[3:], orDone)...):
			}
		}
	}()
	return orDone
}

func signal(after time.Duration) <-chan struct{} {
	c := make(chan struct{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func main() {
	start := time.Now()
	<-OrChan(
		signal(2*time.Hour),
		signal(5*time.Second),
		signal(10*time.Minute),
		signal(8*time.Second),
		signal(7*time.Second),
	)
	fmt.Printf("finishing main after %v", time.Since(start))
}
