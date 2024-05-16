package main

import (
	"log"
	"net"
	"os"

	pb "github.com/lgeorgieff/go-playground/proto/todo/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalln("usage: server [IP_ADDR]")
	}

	addr := args[0]
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen at %s", addr)
	}

	defer func() {
		if err := listener.Close(); err != nil {
			log.Fatalf("failed to closese %v", err)
		}
	}()

	log.Printf("listening at %s\n", addr)

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	pb.RegisterTodoServiceServer(s, NewTodoServer())
	// Usage with grpcui:
	//  * go install github.com/fullstorydev/grpcui/cmd/grpcui@latest
	//  * $GOPATH/bin/grpcui -plaintext 127.0.0.1:50051
	reflection.Register(s)

	defer s.Stop()
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve %v\n", err)
	}
}
