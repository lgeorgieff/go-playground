package main

import (
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalln("usage: client [IP_ADDR]")
	}

	addr := args[0]
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("failed to closese %v", err)
		}
	}()
}
