package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Number interface {
	int8 | int16 | int32 | int64 | int | float32 | float64
}

type AsyncNumber[T Number] struct {
	m *sync.Mutex
	n T
}

func NewAsyncNumber[T Number]() *AsyncNumber[T] {
	return &AsyncNumber[T]{
		m: &sync.Mutex{},
	}
}

func (a *AsyncNumber[T]) Add(n T) {
	a.m.Lock()
	defer a.m.Unlock()

	a.n += n
}

func (a *AsyncNumber[T]) N() T {
	a.m.Lock()
	defer a.m.Unlock()

	n := a.n
	return n
}

func main() {
	const iterations = 1000

	n := NewAsyncNumber[int]()
	wg := &sync.WaitGroup{}
	wg.Add(3 * iterations)
	for i := 0; i < iterations; i++ {
		go func() {
			n.Add(1)
			time.Sleep(10 * time.Millisecond * time.Duration(rand.Intn(10)))
			wg.Done()
		}()
		go func() {
			n.Add(-1)
			time.Sleep(10 * time.Millisecond * time.Duration(rand.Intn(10)))
			wg.Done()
		}()
		go func() {
			fmt.Println(n.N())
			time.Sleep(10 * time.Millisecond * time.Duration(rand.Intn(10)))
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Printf("finished all iterations with n=%d", n.N())
}
