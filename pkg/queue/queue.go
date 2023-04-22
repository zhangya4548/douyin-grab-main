package queue

import (
	"container/list"
	"sync"
)

type QueueSrv struct {
	queue *list.List
	lock  sync.Mutex
}

func NewQueueSrv() *QueueSrv {
	return &QueueSrv{
		queue: list.New(),
	}
}

func (q *QueueSrv) Push(str string) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue.PushBack(str)
}

func (q *QueueSrv) Pop() string {
	if q.queue.Len() == 0 {
		return ""
	}
	e := q.queue.Front()
	q.queue.Remove(e)
	return e.Value.(string)
}

func (q *QueueSrv) GetAll() []string {
	q.lock.Lock()
	defer q.lock.Unlock()
	result := make([]string, 0, q.queue.Len())
	for e := q.queue.Front(); e != nil; e = e.Next() {
		result = append(result, e.Value.(string))
	}
	defer q.Empty()
	return result
}

func (q *QueueSrv) Empty() {
	q.queue.Init()
}
