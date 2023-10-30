package models

import (
	"container/list"
	"sync"
)

type CarQueue struct {
	queue *list.List
	mutex sync.Mutex
}

func NewCarQueue() *CarQueue {
	return &CarQueue{
		queue: list.New(),
	}
}

func (cq *CarQueue) Enqueue(car *Car) {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()
	cq.queue.PushBack(car)
}

func (cq *CarQueue) Dequeue() *Car {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()
	if cq.queue.Len() == 0 {
		return nil
	}
	element := cq.queue.Front()
	cq.queue.Remove(element)
	return element.Value.(*Car)
}

func (cq *CarQueue) First() *Car {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()
	if cq.queue.Len() == 0 {
		return nil
	}
	element := cq.queue.Front()
	return element.Value.(*Car)
}

func (cq *CarQueue) Last() *Car {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()
	if cq.queue.Len() == 0 {
		return nil
	}
	element := cq.queue.Back()
	return element.Value.(*Car)
}

func (cq *CarQueue) Size() int {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()
	return cq.queue.Len()
}
