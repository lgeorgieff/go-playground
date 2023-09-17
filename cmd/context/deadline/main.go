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
			fmt.Println("deadline...")
		}
	}()

	return done
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*2))
	defer cancel()

	<-sayHello(ctx)
	fmt.Println("finished main")
}
