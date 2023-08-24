package main

import (
	"fmt"
	"time"
)

func NewPipeline(done <-chan struct{}, data ...int) <-chan int {
	return NewAddStep(
		done,
		NewMultiplyStep(
			done,
			NewAddStep(
				done,
				NewStream(done, data...),
				3),
			2),
		1)
}

func NewStream(done <-chan struct{}, data ...int) <-chan int {
	stream := make(chan int)
	go func() {
		defer close(stream)

		for _, i := range data {
			time.Sleep(500 * time.Millisecond)
			select {
			case stream <- i:
			case <-done:
				return
			}
		}
	}()
	return stream
}

func NewMultiplyStep(done <-chan struct{}, stream <-chan int, multiplier int) <-chan int {
	result := make(chan int)
	go func() {
		defer close(result)
		for i := range stream {
			time.Sleep(500 * time.Millisecond)
			select {
			case result <- i * multiplier:
			case <-done:
				return
			}
		}
	}()
	return result
}

func NewAddStep(done <-chan struct{}, stream <-chan int, additive int) <-chan int {
	result := make(chan int)
	go func() {
		defer close(result)
		for i := range stream {
			time.Sleep(500 * time.Millisecond)
			select {
			case result <- i + additive:
			case <-done:
				return
			}
		}
	}()
	return result
}

func main() {
	done := make(chan struct{})
	defer close(done)

	for result := range NewPipeline(done, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9) {
		fmt.Println(result)
	}

	fmt.Println("finishing main")
}
