package main

type task struct {
	id      int
	title   string
	content string
}

type taskNode struct {
	next    *taskNode
	current task
}

type taskList struct {
	head, next *taskNode
}

func (list *taskList) appendTask(task task) {
	if list.head == nil {
		list.head = &taskNode{current: task}
		list.next = list.head
	} else {
		list.next.next = &taskNode{current: task}
		list.next = list.next.next
	}
}

func (list *taskNode) getAll() []task {
	var tasks []task

	for t := list.next; t != nil; t = t.next {
		tasks = append(tasks, t.current)
	}

	return tasks
}
