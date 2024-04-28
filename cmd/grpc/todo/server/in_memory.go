package main

import (
	"log"
	"time"

	pb "github.com/lgeorgieff/go-playground/proto/todo/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type inMemmoryDB struct {
	tasks  map[uint64]*pb.Task
	nextID uint64
}

func (db *inMemmoryDB) addTask(description string, dueDate time.Time) (uint64, error) {
	task := &pb.Task{
		Id:          db.nextID,
		Description: description,
		Done:        false,
		DueDate:     timestamppb.New(dueDate),
	}
	db.nextID++
	db.tasks[task.Id] = task

	log.Printf("Stored task description=\"%s\", dueDate=%v, id=%d\n", description, dueDate, task.Id)

	return task.Id, nil
}

func NewDB() db {
	return &inMemmoryDB{
		tasks: make(map[uint64]*pb.Task),
	}
}
