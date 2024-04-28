package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	pb "github.com/lgeorgieff/go-playground/proto/todo/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	c := pb.NewTodoServiceClient(conn)

	for i := 0; i < 10; i++ {
		dueData := time.Now().Add(time.Duration(rand.Intn(60)) * time.Second)
		description := fmt.Sprintf("This is task no. %d", i)
		id, err := addTask(c, description, dueData)
		log.Default().Printf("Adding task \"%s\" with result.id=%d, err=%v\n", description, id, err)
	}
}

func addTask(c pb.TodoServiceClient, description string, dueDate time.Time) (uint64, error) {
	req := &pb.AddTaskRequest{
		Description: description,
		DueDate:     timestamppb.New(dueDate),
	}
	res, err := c.AddTask(context.Background(), req)
	if err != nil {
		return 0, err
	}
	return res.Id, nil
}
