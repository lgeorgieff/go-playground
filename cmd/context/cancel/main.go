package main

import (
	"context"
	"fmt"
	"time"
)

func sayHello(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)
		select {
		case <-time.After(time.Second * 3):
			fmt.Println("Hello, World!")
		case <-ctx.Done():
			fmt.Println("cancelled...")
		}
	}()

	return done
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()

	<-sayHello(ctx)
	fmt.Println("finished main")
}
