package models

import "sync"

// Queue es una estructura de datos de cola
type Queue struct {
	cars     []*Car
	mutex    sync.Mutex
	notEmpty *sync.Cond
}

func NewQueue() *Queue {
	q := &Queue{}
	q.notEmpty = sync.NewCond(&q.mutex)
	return q
}

func (q *Queue) Enqueue(car *Car) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.cars = append(q.cars, car)
	q.notEmpty.Signal() // Notificar a las goroutines en espera
}

func (q *Queue) Dequeue() *Car {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for len(q.cars) == 0 {
		q.notEmpty.Wait() // Esperar si la cola está vacía
	}

	item := q.cars[0]
	q.cars = q.cars[1:]
	return item
}

func (q *Queue) First() *Car {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.cars) == 0 {
		return nil
	}
	return q.cars[0]
}

func (q *Queue) Last() *Car {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.cars) == 0 {
		return nil
	}
	return q.cars[len(q.cars)-1]
}

func (q *Queue) Size() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return len(q.cars)
}
