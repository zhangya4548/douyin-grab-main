// @Time : 2023/3/30 7:16 PM
// @Author : zhangguangqiang
// @File : queue
// @Software: GoLand

package queue

import "container/list"

type QueueSrv struct {
	queue *list.List
}

func NewQueueSrv() *QueueSrv {
	return &QueueSrv{
		queue: list.New(),
	}
}

// Push 将json字符串写入队列尾部
func (q QueueSrv) Push(jsonStr string) {
	q.queue.PushBack(jsonStr)
}

// Pop 从队列头部获取json字符串
func (q QueueSrv) Pop() string {
	if q.queue.Len() == 0 {
		return ""
	}
	e := q.queue.Front()
	q.queue.Remove(e)
	return e.Value.(string)
}
