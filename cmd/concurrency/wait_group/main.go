package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	const iterations = 10

	wg := &sync.WaitGroup{}
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func(ii int) {
			time.Sleep(250 * time.Millisecond * time.Duration(rand.Intn(10)))
			fmt.Printf("iteration %d / %d\n", ii, iterations)
			wg.Done()
		}(i)
	}

	fmt.Printf("waiting for %d iterations to finish", iterations)
	wg.Wait()
}
