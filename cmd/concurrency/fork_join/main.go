package main

import (
	"fmt"
	"sync"
)

func fork(joiner *sync.WaitGroup) {
	defer joiner.Done()

	fmt.Println("forked ...")
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go fork(wg)
	fmt.Println("main ...")

	wg.Wait()
}
