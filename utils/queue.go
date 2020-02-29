package utils

import (
	"container/list"
	"sync"
)

// Queue ...
type Queue struct {
	sync.Mutex
	innerList *list.List
	cond      *sync.Cond
}

// Size ..
func (queue *Queue) Size() int {
	return queue.innerList.Len()
}

func (queue *Queue) Clear() {
	queue.innerList.Init()
}

// Push ..
func (queue *Queue) Push(value interface{}) {
	queue.Lock()
	queue.innerList.PushBack(value)
	queue.Unlock()
	queue.cond.Broadcast()
}

func (queue *Queue) _innerPop() (value interface{}) {
	el := queue.innerList.Front()
	if el != nil {
		value = el.Value
		queue.innerList.Remove(el)
	}
	return value
}

func (queue *Queue) Pop() (value interface{}) {
	queue.Lock()
	defer queue.Unlock()
	
	value = queue._innerPop()
	for ; value == nil; value = queue._innerPop() {
		queue.cond.Wait()
	}
	
	return value
}

func NewQueue() *Queue {
	queue := &Queue{innerList: list.New()}
	queue.cond = sync.NewCond(queue)
	return queue
}
