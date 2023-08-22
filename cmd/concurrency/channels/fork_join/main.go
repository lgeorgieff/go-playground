package main

import "fmt"

func goroutine(done chan<- struct{}) {
	fmt.Println("Hello, World!")
	// In case we do not need the channel for anything else, we may close it to signal that the go routine run
	// otherwise we simply write an item to it
	done <- struct{}{}
	// close(done)
}

func main() {
	barrier := make(chan struct{})

	go goroutine(barrier)
	<-barrier
}
