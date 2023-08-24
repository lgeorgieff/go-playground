package main

import (
	"fmt"
	"sync"
)

func main() {
	const iterations = 10000

	wg := sync.WaitGroup{}
	wg.Add(1)

	done := make(chan struct{})
	close(done)
	writeChan := make(chan int, iterations)
	var doneCounter, writeCounter, defaultCounter uint
	go func() {
		defer wg.Done()

		for i := 0; i < iterations; i++ {
			select {
			case <-done:
				doneCounter++
			case writeChan <- i:
				writeCounter++
			default:
				defaultCounter++
			}
		}
	}()

	wg.Wait()
	fmt.Println("doneCounter:", doneCounter)
	fmt.Println("writeCounter:", writeCounter)
	fmt.Println("defaultCounter:", defaultCounter)
}
