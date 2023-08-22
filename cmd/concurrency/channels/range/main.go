package main

import (
	"fmt"
	"math/rand"
	"time"
)

func worker(iterations int, stream chan<- int) {
	defer close(stream)

	for i := 0; i < iterations; i++ {
		time.Sleep(50 * time.Millisecond * time.Duration(rand.Intn(10)))
		stream <- i
	}
}

func main() {
	const (
		workerCount = 10
		iterations  = 50
	)
	stream := make(chan int)

	go worker(iterations, stream)

	// we can range over a channel, once it will be closed, the loop exits
	for i := range stream {
		fmt.Println(i)
	}
}
