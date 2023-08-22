package main

import "fmt"

func main() {
	// Create a bidrectional channel that allows writing and reading
	bidirectionalChannel := make(chan string)

	go func(writeChannel chan<- string) {
		// Use writable channel for writing a message into the channel
		// This statement will block until the channel will be empty
		writeChannel <- "Hello, World!"
	}(bidirectionalChannel)

	var readChannel <-chan string = bidirectionalChannel
	// Use readable channel for reading a message from the channel
	// This statement will block until the channel will be filled with an item
	msg := <-readChannel
	fmt.Println(msg)
}
