package main

import (
	"fmt"
	"sync"
)

type MyStruct struct {
	once  sync.Once
	state int
}

func (m *MyStruct) start() {
	m.once.Do(func() {
		m.state++
	})
}

func (m *MyStruct) Start() {
	m.start()
}

func (m *MyStruct) Stop() {
	fmt.Println("state:", m.state)
}

func main() {
	const iterations = 100

	m := MyStruct{}
	wg := sync.WaitGroup{}
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func() {
			m.Start()
			wg.Done()
		}()
	}
	wg.Wait()
	m.Stop()
}
