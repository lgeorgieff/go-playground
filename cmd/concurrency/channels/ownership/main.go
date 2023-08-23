package main

import (
	"fmt"
	"time"
)

// the resource/channel owner
func producer() <-chan int {
	const iterations = 6
	// since we know how many values we will generate, we optimize and create a buffered channel
	resultStream := make(chan int, iterations-1)
	go func() {
		// once we finished work, we close the channel within the owner scope
		defer close(resultStream)
		for i := 0; i < iterations; i++ {
			resultStream <- i
			fmt.Println("wrote:", i)
		}
	}()

	// we return a read-only channel, since the consumer only needs to read from this channel
	return resultStream
}

func consumer(stream <-chan int) {
	// we are only concerned about blocking/waiting and reading values from the channel
	for result := range stream {
		fmt.Println("read:", result)
	}
}

func main() {
	resultStream := producer()
	time.Sleep(500 * time.Millisecond)
	consumer(resultStream)

	fmt.Println("finishing main")
}
