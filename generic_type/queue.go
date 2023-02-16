package generic_type

import (
	"container/list"
)

type Queue struct {
	l *list.List
}

func NewQueue() *Queue {
	return &Queue{
		l: list.New(),
	}
}

func (q *Queue) Len() int {
	return q.l.Len()
}

func (q *Queue) IsEmpty() bool {
	return q.l.Len() == 0
}

func (q *Queue) Push(item interface{}) {
	q.l.PushBack(item)
}

func (q *Queue) Pop() interface{} {
	return q.l.Remove(q.l.Back())
}