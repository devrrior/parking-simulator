package models

import (
	"image/color"
	"sync"
	"time"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/scene"
)

const entranceSpotX = 355.00

type Car struct {
	area   floatgeom.Rect2
	entity *entities.Entity
}

func NewCar(ctx *scene.Context) *Car {
	area := floatgeom.NewRect2(300, 480, 320, 460)
	entity := entities.New(ctx, entities.WithRect(area), entities.WithColor(color.RGBA{255, 0, 0, 255}))

	return &Car{
		area:   area,
		entity: entity,
	}
}

func (c *Car) Enqueue(queueCars *Queue) {

	if queueCars.Size() == 0 {
		for c.entity.Y() > 45 {

			c.entity.ShiftY(-1)

			time.Sleep(20 * time.Millisecond)
		}

		return
	}

	lastCar := queueCars.Last()

	for c.entity.Y() > 45 {

		if (c.entity.Y() - 10) < lastCar.entity.Y() {
			c.entity.ShiftY(-1)
		}

		time.Sleep(20 * time.Millisecond)
	}
	return
}

func (c *Car) JoinDoor(semaphore *Semaphore, doorM *sync.Mutex) {
	semaphore.Wait()

	doorM.Lock()
	for {
		if c.entity.X() < entranceSpotX {
			c.entity.ShiftX(1)
		}
		if c.entity.X() == entranceSpotX {
			break
		}

		time.Sleep(20 * time.Millisecond)
	}
	doorM.Unlock()
}

func (c *Car) ExitDoor(semaphore *Semaphore, doorM *sync.Mutex) {
	doorM.Lock()
	for {
		if c.entity.X() > 0 {
			c.entity.ShiftX(-1)
		}
		if c.entity.X() == 0 {
			break
		}

		time.Sleep(20 * time.Millisecond)
	}
	doorM.Unlock()

	semaphore.Signal()
}

func (c *Car) Park(spot *ParkingSpot) {
	for index := 0; index < len(*spot.GetDirectionsForParking()); index++ {
		directions := *spot.GetDirectionsForParking()
		if directions[index].Direction == "right" {
			for {
				if c.entity.X() < directions[index].Point {
					c.entity.ShiftX(1)
				}
				if c.entity.X() == directions[index].Point {
					break
				}

				time.Sleep(20 * time.Millisecond)
			}
		} else if directions[index].Direction == "down" {
			for {
				if c.entity.Y() < directions[index].Point {
					c.entity.ShiftY(1)
				}
				if c.entity.Y() == directions[index].Point {
					break
				}

				time.Sleep(20 * time.Millisecond)
			}
		}
	}
}

func (c *Car) Leave(spot *ParkingSpot) {
	for index := 0; index < len(*spot.GetDirectionsForLeaving()); index++ {
		directions := *spot.GetDirectionsForLeaving()
		if directions[index].Direction == "left" {

			for {
				if c.entity.X() > directions[index].Point {
					c.entity.ShiftX(-1)
				}
				if c.entity.X() == directions[index].Point {
					break
				}

				time.Sleep(20 * time.Millisecond)
			}
		} else if directions[index].Direction == "right" {
			for {
				if c.entity.X() < directions[index].Point {
					c.entity.ShiftX(1)
				}
				if c.entity.X() == directions[index].Point {
					break
				}

				time.Sleep(20 * time.Millisecond)
			}
		} else if directions[index].Direction == "up" {
			for {
				if c.entity.Y() > directions[index].Point {
					c.entity.ShiftY(-1)
				}
				if c.entity.Y() == directions[index].Point {
					break
				}

				time.Sleep(20 * time.Millisecond)
			}
		} else if directions[index].Direction == "down" {
			for {
				if c.entity.Y() < directions[index].Point {
					c.entity.ShiftY(1)
				}
				if c.entity.Y() == directions[index].Point {
					break
				}

				time.Sleep(20 * time.Millisecond)
			}
		}
	}
}

func (c *Car) Remove() {
	c.entity.Destroy()
}

func (c *Car) GetEntityX() float64 {
	return c.entity.X()
}

func (c *Car) GetEntityY() float64 {
	return c.entity.Y()
}

func (c *Car) EntityShiftX(x float64) {
	c.entity.ShiftX(x)
}

func (c *Car) EntityShiftY(y float64) {
	c.entity.ShiftY(y)
}