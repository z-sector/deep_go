package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	heap *Heap[Task, int]
}

func NewScheduler() Scheduler {
	less := func(a, b Task) bool {
		return a.Priority > b.Priority
	}

	getIdentifier := func(task Task) int {
		return task.Identifier
	}
	return Scheduler{
		heap: NewHeap[Task, int](less, getIdentifier),
	}
}

func (s *Scheduler) AddTask(task Task) {
	s.heap.Push(task)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	task, ok := s.heap.GetByIdentifier(taskID)
	if !ok {
		return
	}

	task.Priority = newPriority
	s.heap.Change(taskID, func(task Task) Task {
		task.Priority = newPriority
		return task
	})
}

func (s *Scheduler) GetTask() Task {
	return s.heap.Pop()
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	task1.Priority = 100
	assert.Equal(t, task1, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
