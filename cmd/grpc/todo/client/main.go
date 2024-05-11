package main

import (
	"context"
	"fmt"
	"io"
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

	printTasks(c)

	updateTasks(c, []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})

	deleteTasks(c, []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
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

func printTasks(c pb.TodoServiceClient) {
	stream, err := c.ListTasks(context.Background(), &pb.ListTasksRequest{})
	if err != nil {
		log.Fatalf("unexpected error when fetching task stream: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("unexpected error when reading task from stream: %v", err)
		}
		log.Printf("receveid task response: %v\n", res)
	}
}

func updateTasks(c pb.TodoServiceClient, taskIDs []uint64) {
	stream, err := c.UpdateTasks(context.Background())
	if err != nil {
		log.Fatalf("unexpected error when creating an UpdateTask stream: %v", err)
	}
	now := time.Now().Add(1 * time.Minute)
	for i, id := range taskIDs {
		task := &pb.Task{
			Id:          id,
			Description: fmt.Sprintf("Updated task %d", i),
			DueDate:     timestamppb.New(now),
			Done:        true,
		}
		if err := stream.Send(&pb.UpdateTasksRequest{Task: task}); err != nil {
			log.Fatalf("unexpected error when sending task update for id %d: %v", id, err)
		}
	}
	if _, err := stream.CloseAndRecv(); err != nil {
		log.Fatalf("unexpected error when closing UpdateTask stream: %v", err)
	}
}

func deleteTasks(c pb.TodoServiceClient, taskIDs []uint64) {
	stream, err := c.DeleteTasks(context.Background())
	if err != nil {
		log.Fatalf("unexpected error when creating a DeleteTask stream: %v", err)
	}

	waitc := make(chan struct{})
	go func() {
		for {
			_, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				break
			}
			if err != nil {
				log.Fatalf("unexpected error when receiving a DeleteTask response: %v", err)
			}

			log.Println("deleted Tasks")
		}
	}()

	for _, id := range taskIDs {
		req := &pb.DeleteTasksRequest{Id: id}
		if err := stream.Send(req); err != nil {
			log.Fatalf("unexpected error when sending DeleteTask for id %d: %v", id, err)
		}
	}
	if err := stream.CloseSend(); err != nil {
		log.Fatalf("unexpected error when closing DeleteTask stream: %v", err)
	}
	<-waitc
}
