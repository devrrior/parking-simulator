package models

import "sync"

type Parking struct {
	spots     []*ParkingSpot
	doorM     *sync.Mutex
	semaphore *Semaphore
	queueCars *Queue
}

func NewParking(spots []*ParkingSpot) *Parking {
	semaphore := NewSemaphore(len(spots))
	queue := NewQueue()

	return &Parking{
		spots:     spots,
		doorM:     &sync.Mutex{},
		semaphore: semaphore,
		queueCars: queue,
	}
}

func (p *Parking) GetSpots() []*ParkingSpot {
	return p.spots
}

func (p *Parking) GetDoorM() *sync.Mutex {
	return p.doorM
}

func (p *Parking) GetSemaphore() *Semaphore {
	return p.semaphore
}

func (p *Parking) GetParkingSpotAvailable() *ParkingSpot {
	for _, spot := range p.spots {
		if spot.GetIsAvailable() {
			spot.SetIsAvailable(false)
			return spot
		}
	}

	return nil
}

func (p *Parking) GetQueueCars() *Queue {
	return p.queueCars
}
