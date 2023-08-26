package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func Ints(done <-chan struct{}) <-chan int {
	stream := make(chan int)

	go func() {
		defer close(stream)
		for {
			select {
			case <-done:
				return
			case stream <- rand.Int():
			}
		}
	}()

	return stream
}

func TakeN[T any](done <-chan struct{}, stream <-chan T, n uint) <-chan T {
	result := make(chan T)

	go func() {
		defer close(result)
		for i := uint(0); i < n; i++ {
			select {
			case <-done:
				return
			case result <- <-stream:
			}
		}
	}()

	return result
}

func Process(done <-chan struct{}, stream <-chan int) <-chan int {
	result := make(chan int)

	go func() {
		defer close(result)
		for v := range stream {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
			select {
			case <-done:
				return
			case result <- v:
			}
		}
	}()

	return result
}

func Print[T any](done <-chan struct{}, stream <-chan T) <-chan struct{} {
	localDone := make(chan struct{})
	go func() {
		defer close(localDone)
		for v := range stream {
			select {
			case <-done:
				return
			default:
				fmt.Println(v)
			}
		}
	}()
	return localDone
}

func FanOut[T any](done <-chan struct{}, stream <-chan T, fn func(done <-chan struct{}, stream <-chan T) <-chan T) []<-chan T {
	numFns := runtime.NumCPU()
	fns := make([]<-chan T, numFns)
	for i := 0; i < numFns; i++ {
		fns[i] = fn(done, stream)
	}
	return fns
}

func FanIn[T any](done <-chan struct{}, streams ...<-chan T) <-chan T {
	result := make(chan T)
	wg := sync.WaitGroup{}
	wg.Add(len(streams))

	for _, stream := range streams {
		go func(s <-chan T) {
			defer wg.Done()
			for v := range s {
				select {
				case <-done:
					return
				case result <- v:
				}
			}
		}(stream)
	}

	go func() {
		wg.Wait()
		close(result)
	}()
	return result
}

func NewPipeline(done <-chan struct{}, iterations uint) <-chan struct{} {
	stream := TakeN(done, Ints(done), iterations)
	slowStreams := FanOut(done, stream, Process)
	stream = FanIn(done, slowStreams...)
	return Print(done, stream)
}

func main() {
	const iterations = 20

	done := make(chan struct{})
	defer close(done)

	<-NewPipeline(done, iterations)
	fmt.Println("\nfinishing main")
}
