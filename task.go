package main

import "math/rand"

type Task struct {
	id      int
	title   string
	content string
}

func createTask(title, content string) *Task {
	return &Task{id: rand.Intn(100), title: title, content: content}
}
