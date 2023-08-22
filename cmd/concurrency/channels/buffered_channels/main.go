package main

import (
	"fmt"
	"time"
)

func main() {
	bufferedStream := make(chan int, 4)
	go func() {
		defer close(bufferedStream)

		for i := 0; i < cap(bufferedStream)*2; i++ {
			time.Sleep(10 * time.Millisecond)
			fmt.Printf("chan <- %d, [%d/%d]\n", i, len(bufferedStream), cap(bufferedStream))
			bufferedStream <- i
		}
	}()

	for i := range bufferedStream {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("%d <- chan, [%d/%d]\n", i, len(bufferedStream), cap(bufferedStream))
	}
}
