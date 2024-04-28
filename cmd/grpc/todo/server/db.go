package main

import "time"

type db interface {
	addTask(descroption string, dueDate time.Time) (uint64, error)
}
