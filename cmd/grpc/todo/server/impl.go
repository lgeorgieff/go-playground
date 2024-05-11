package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/lgeorgieff/go-playground/proto/todo/v1"
)

var _ pb.TodoServiceServer = (*server)(nil)

type server struct {
	d db
	pb.UnimplementedTodoServiceServer
}

func NewTodoServer() pb.TodoServiceServer {
	return &server{
		d: NewDB(),
	}
}

func (s *server) AddTask(_ context.Context, in *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
	id, err := s.d.addTask(in.Description, in.DueDate.AsTime())
	if err != nil {
		return nil, err
	}
	return &pb.AddTaskResponse{Id: id}, nil
}

func (s *server) ListTasks(req *pb.ListTasksRequest, stream pb.TodoService_ListTasksServer) error {
	now := time.Now().UTC()
	handler := func(task *pb.Task) error {
		res := &pb.ListTasksResponse{
			Task:    task,
			Overdue: !task.Done && now.After(task.DueDate.AsTime()),
		}
		return stream.Send(res)
	}

	return s.d.getTasks(handler)
}

func (s *server) UpdateTasks(stream pb.TodoService_UpdateTasksServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.UpdateTasksResponse{})
		}
		if err != nil {
			return err
		}
		if err := s.d.updateTask(req.Task); err != nil {
			log.Printf("failed to update task with id %d in DB: %v\n", req.Task.Id, err)
		}
	}
}
