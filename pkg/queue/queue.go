package queue

import (
	"container/list"
	"fmt"
	"sync"
)

type QueueSrv struct {
	queue *list.List
	lock  sync.Mutex // 加锁，保证并发安全
}

func NewQueueSrv() *QueueSrv {
	return &QueueSrv{
		queue: list.New(),
	}
}

// Push 将字符串写入队列尾部
func (q *QueueSrv) Push(str string) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue.PushBack(str)
}

// Pop 从队列头部获取字符串
func (q *QueueSrv) Pop() string {
	if q.queue.Len() == 0 {
		return ""
	}
	e := q.queue.Front()
	q.queue.Remove(e)
	return e.Value.(string)
}

// GetAll 获取队列中所有元素的值，并清空队列
func (q *QueueSrv) GetAll() []string {
	q.lock.Lock() // 加锁，保证并发安全
	defer q.lock.Unlock()
	fmt.Println("进来了======", q.queue.Len())
	result := make([]string, 0, q.queue.Len())
	for e := q.queue.Front(); e != nil; e = e.Next() {
		result = append(result, e.Value.(string))
	}
	defer q.Empty()
	fmt.Println("进来了1======", q.queue.Len())
	return result
}

// 清空队列
func (q *QueueSrv) Empty() {
	q.queue.Init()
}
