package main

import (
	"time"

	pb "github.com/lgeorgieff/go-playground/proto/todo/v1"
)

type db interface {
	addTask(descroption string, dueDate time.Time) (uint64, error)
	getTasks(f func(task *pb.Task) error) error
	updateTask(task *pb.Task) error
	deleteTask(id uint64) error
}
