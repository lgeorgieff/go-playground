package main

import (
	"fmt"

	"github.com/lgeorgieff/go-playground/internal/domain/channel"
)

func Tee[T any](stream <-chan T, done <-chan struct{}) (<-chan T, <-chan T) {
	out1 := make(chan T)
	out2 := make(chan T)

	go func() {
		for v := range channel.OrDone(stream, done) {
			out1, out2 := out1, out2
			for i := 0; i < 2; i++ {
				select {
				case out1 <- v:
					out1 = nil
				case out2 <- v:
					out2 = nil
				case <-done:
					return
				}
			}
		}
	}()

	return out1, out2
}

func main() {
	const size = 10

	done := make(chan struct{})
	stream := channel.Source(size, done, channel.DefaultIntGenerator)
	t1, t2 := Tee(stream, done)

	for i := 0; i < size; i++ {
		fmt.Printf("t1: %d, t2: %d\n", <-t1, <-t2)
	}
}
