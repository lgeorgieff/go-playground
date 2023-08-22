package main

import "fmt"

func goroutine(done chan<- struct{}) {
	fmt.Println("Hello, World!")
	close(done)
}

func main() {
	done := make(chan struct{})

	go goroutine(done)

	// This checks whether the first value was inserted into the channel
	// or whether it is a default value that resulted from a closed channel
	_, ok := <-done
	fmt.Println(ok)
}
