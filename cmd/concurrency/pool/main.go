package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var resourceCounter atomic.Uint32 = atomic.Uint32{}

type expensiveResource struct{}

func (e *expensiveResource) InstanceCount() uint32 {
	return resourceCounter.Load()
}

func newExpensiveResource() *expensiveResource {
	resourceCounter.Add(1)
	return &expensiveResource{}
}

func main() {
	const iterations = 100

	pool := sync.Pool{
		// Any item stored in the Pool may be removed automatically at any time without
		// notification. If the Pool holds the only reference when this happens, the
		// item might be deallocated.

		// New optionally specifies a function to generate
		// a value when Get would otherwise return nil.
		New: func() any { return newExpensiveResource() },
	}
	wgPut := sync.WaitGroup{}
	wgPut.Add(iterations / 10)
	for i := 0; i < iterations/10; i++ {
		go func() {
			pool.Put(newExpensiveResource())
			wgPut.Done()
		}()
	}
	wgPut.Wait()
	fmt.Println("resourceCounter:", resourceCounter.Load())

	wgGet := sync.WaitGroup{}
	wgGet.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func() {
			time.Sleep(250 * time.Millisecond * time.Duration(rand.Intn(50)))
			res := pool.Get()
			fmt.Println("expensiveResource:", res)
			pool.Put(res)
			wgGet.Done()
		}()
	}
	wgGet.Wait()

	fmt.Println("resourceCounter:", resourceCounter.Load())
}
