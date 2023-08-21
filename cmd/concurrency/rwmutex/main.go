package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Reader[K comparable, V any] func(size int, k K, v V)

type Dictionary[K comparable, V any] struct {
	data map[K]V
	m    *sync.RWMutex
}

func (d *Dictionary[K, V]) Write(k K, v V) bool {
	d.m.Lock()
	defer d.m.Unlock()

	isNewKey := true
	if _, ok := d.data[k]; ok {
		isNewKey = false
	}
	d.data[k] = v
	return isNewKey
}

func (d *Dictionary[K, V]) ReadAll(r Reader[K, V]) {
	d.m.RLock()
	defer d.m.RUnlock()

	for k, v := range d.data {
		r(len(d.data), k, v)
	}
}

func NewDictionary[K comparable, V any]() *Dictionary[K, V] {
	return &Dictionary[K, V]{
		data: make(map[K]V),
		m:    &sync.RWMutex{},
	}
}

func main() {
	const (
		iterationsRead  = 100
		iterationsWrite = 10
	)

	dict := NewDictionary[int, float32]()
	wgRead := sync.WaitGroup{}
	wgRead.Add(iterationsRead)
	for i := 0; i < iterationsRead; i++ {
		go func(i int) {
			time.Sleep(10 * time.Millisecond * time.Duration(rand.Intn(10)))
			dict.Write(i, float32(i))
			wgRead.Done()
		}(i)
	}

	wgWrite := sync.WaitGroup{}
	wgWrite.Add(iterationsWrite)
	for i := 0; i < iterationsWrite; i++ {
		go func(i int) {
			time.Sleep(10 * time.Millisecond * time.Duration(rand.Intn(10)))
			dict.ReadAll(func(size, k int, v float32) {
				fmt.Printf("[%d]: size=%d, %d=%f\n", i, size, k, v)
			})
			wgWrite.Done()
		}(i)
	}

	wgRead.Wait()
	wgWrite.Wait()
}
