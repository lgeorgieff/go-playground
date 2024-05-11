package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/lgeorgieff/go-playground/proto/todo/v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if len(in.Description) == 0 {
		return nil, status.Error(codes.InvalidArgument, "expected task description, got an empty string")
	}
	if in.DueDate.AsTime().Before(time.Now().UTC()) {
		return nil, status.Error(codes.InvalidArgument, "expected a task due_date that is in the future")
	}

	id, err := s.d.addTask(in.Description, in.DueDate.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unexpected error: %s", err.Error())
	}
	return &pb.AddTaskResponse{Id: id}, nil
}

func (s *server) ListTasks(req *pb.ListTasksRequest, stream pb.TodoService_ListTasksServer) error {
	now := time.Now().UTC()
	handler := func(task *pb.Task) error {
		pb.Filter(task, req.Mask)

		res := &pb.ListTasksResponse{
			Task:    task,
			Overdue: task.DueDate != nil && !task.Done && now.After(task.DueDate.AsTime()),
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

		task := &pb.Task{
			Id:          req.Id,
			Description: req.Description,
			Done:        req.Done,
			DueDate:     req.DueDate,
		}
		if err := s.d.updateTask(task); err != nil {
			log.Printf("failed to update task with id %d in DB: %v\n", req.Id, err)
		}
	}
}

func (s *server) DeleteTasks(stream pb.TodoService_DeleteTasksServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := s.d.deleteTask(req.Id); err != nil {
			return errors.Wrapf(err, "failed to delete task with ID %d from DB", req.Id)
		}
		if err := stream.Send(&pb.DeleteTasksResponse{}); err != nil {
			return errors.Wrapf(err, "failed to send DeleteTasksResponse for ID %d", req.Id)
		}
	}
}
