package main

import (
	"fmt"
	"math/rand"
	"time"
)

func worker(done chan struct{}) {
	time.Sleep(time.Second * time.Duration(rand.Intn(6)))
	close(done)
}

func main() {
	done := make(chan struct{})
	timeout := time.After(time.Second * 3)

	go worker(done)
loop:
	for {
		select {
		case <-done:
			fmt.Println("done")
			break loop
		case <-timeout:
			fmt.Println("timeout")
			break loop
		default:
			fmt.Println("default")
			time.Sleep(100 * time.Millisecond)
		}
	}
	fmt.Println("finished main")
}
