package main

import (
	"context"

	pb "github.com/lgeorgieff/go-playground/proto/todo/v1"
)

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
