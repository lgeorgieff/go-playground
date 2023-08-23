package main

import (
	"fmt"
	"sync"
	"time"
)

func goroutine(id int, barrier <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println(id, "waiting before barrier")
	<-barrier
	fmt.Println(id, "passed barrier")
}

func main() {
	const goroutineCount = 10

	barrier := make(chan struct{})
	wg := &sync.WaitGroup{}
	wg.Add(goroutineCount)
	for i := 0; i < goroutineCount; i++ {
		go goroutine(i, barrier, wg)
	}

	time.Sleep(1 * time.Second)
	close(barrier)
	wg.Wait()
}
