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
		case <-time.After(3 * time.Second):
			fmt.Println("Hello, World!")
		case <-ctx.Done():
			fmt.Println("timeout...")
		}
	}()

	return done
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	<-sayHello(ctx)
	fmt.Println("finished main")
}
