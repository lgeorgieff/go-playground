package main

import (
	"log"
	"time"

	pb "github.com/lgeorgieff/go-playground/proto/todo/v1"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ db = (*inMemmoryDB)(nil)

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

func (db *inMemmoryDB) getTasks(f func(task *pb.Task) error) error {
	for _, task := range db.tasks {
		if err := f(task); err != nil {
			return err
		}
	}
	return nil
}

func (db *inMemmoryDB) updateTask(task *pb.Task) error {
	if dbTask, ok := db.tasks[task.Id]; ok {
		dbTask.Description = task.Description
		dbTask.Done = task.Done
		dbTask.DueDate = task.DueDate

		log.Printf("Updated task description=\"%s\", dueDate=%v, done=%t, id=%d\n", task.Description, task.DueDate, task.Done, task.Id)
		return nil
	}
	return errors.Errorf("task with id %d not found in DB", task.Id)
}

func NewDB() db {
	return &inMemmoryDB{
		tasks: make(map[uint64]*pb.Task),
	}
}
