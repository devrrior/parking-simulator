package models

import (
	"sync"
)

type Parking struct {
	spotsMutex     sync.Mutex
	availableSpots *sync.Cond
	spots          []*ParkingSpot
	queueCars      *CarQueue
}

func NewParking(spots []*ParkingSpot) *Parking {
	queue := NewCarQueue()
	cond := sync.NewCond(&sync.Mutex{})

	return &Parking{
		spots:          spots,
		availableSpots: cond,
		queueCars:      queue,
	}
}

func (p *Parking) GetSpots() []*ParkingSpot {
	p.spotsMutex.Lock()
	defer p.spotsMutex.Unlock()
	return p.spots
}

func (p *Parking) GetParkingSpotAvailable() *ParkingSpot {
	p.spotsMutex.Lock()
	defer p.spotsMutex.Unlock()

	for {
		for _, spot := range p.spots {
			if spot.GetIsAvailable() {
				spot.SetIsAvailable(false)
				return spot
			}
		}

		p.availableSpots.Wait()
	}
}

func (p *Parking) ReleaseParkingSpot(spot *ParkingSpot) {
	p.spotsMutex.Lock()
	defer p.spotsMutex.Unlock()
	spot.SetIsAvailable(true)
	p.availableSpots.Signal() // Notificar a una goroutine en espera
}

func (p *Parking) GetQueueCars() *CarQueue {
	return p.queueCars
}
