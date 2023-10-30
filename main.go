package main

import (
	"image/color"
	"math/rand"
	"parking-concurrency/models"
	"sync"
	"time"

	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/scene"
)

var (
	spots = []*models.ParkingSpot{
		// first row
		models.NewParkingSpot(380, 70, 410, 100, 1, 1),
		models.NewParkingSpot(425, 70, 455, 100, 1, 2),
		models.NewParkingSpot(470, 70, 500, 100, 1, 3),
		models.NewParkingSpot(515, 70, 545, 100, 1, 4),
		models.NewParkingSpot(560, 70, 590, 100, 1, 5),

		// second row
		models.NewParkingSpot(380, 160, 410, 190, 2, 6),
		models.NewParkingSpot(425, 160, 455, 190, 2, 7),
		models.NewParkingSpot(470, 160, 500, 190, 2, 8),
		models.NewParkingSpot(515, 160, 545, 190, 2, 9),
		models.NewParkingSpot(560, 160, 590, 190, 2, 10),

		// third row
		models.NewParkingSpot(380, 250, 410, 280, 3, 11),
		models.NewParkingSpot(425, 250, 455, 280, 3, 12),
		models.NewParkingSpot(470, 250, 500, 280, 3, 13),
		models.NewParkingSpot(515, 250, 545, 280, 3, 14),
		models.NewParkingSpot(560, 250, 590, 280, 3, 15),

		// fourth row
		models.NewParkingSpot(380, 340, 410, 370, 4, 16),
		models.NewParkingSpot(425, 340, 455, 370, 4, 17),
		models.NewParkingSpot(470, 340, 500, 370, 4, 18),
		models.NewParkingSpot(515, 340, 545, 370, 4, 19),
		models.NewParkingSpot(560, 340, 590, 370, 4, 20),
	}
	parking   = models.NewParking(spots)
	queueCars = parking.GetQueueCars()
	doorMutex sync.Mutex
)

func main() {
	isFirstTime := true
	wg := sync.WaitGroup{}

	_ = oak.AddScene("parkingScene", scene.Scene{
		Start: func(ctx *scene.Context) {
			_ = ctx.Window.SetBorderless(true)
			setUpScene(ctx)

			event.GlobalBind(ctx, event.Enter, func(enterPayload event.EnterPayload) event.Response {
				if !isFirstTime {
					return 0
				}

				isFirstTime = false

				for i := 0; i < 50; i++ {
					wg.Add(1)
					go carCycle(ctx, &wg)

					time.Sleep(time.Millisecond * time.Duration(getRandomNumber(1000, 2000)))
				}

				return 0
			})
		},
	})

	_ = oak.Init("parkingScene")
	wg.Wait()
}

func setUpScene(ctx *scene.Context) {

	parkingArea := floatgeom.NewRect2(350, 10, 630, 400)
	entities.New(ctx, entities.WithRect(parkingArea), entities.WithColor(color.RGBA{86, 101, 115, 255}))

	parkingDoor := floatgeom.NewRect2(340, 10, 350, 70)
	entities.New(ctx, entities.WithRect(parkingDoor), entities.WithColor(color.RGBA{255, 255, 255, 255}))

	for _, spot := range spots {
		entities.New(ctx, entities.WithRect(*spot.GetArea()), entities.WithColor(color.RGBA{212, 172, 13, 255}))
	}
}

func carCycle(ctx *scene.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	car := models.NewCar(ctx)

	car.Enqueue()

	queueCars.Enqueue(car)

	doorMutex.Lock()

	car.JoinDoor()

	doorMutex.Unlock()

	queueCars.Dequeue()

	spotAvailable := parking.GetParkingSpotAvailable()

	car.Park(spotAvailable)

	time.Sleep(time.Millisecond * time.Duration(getRandomNumber(1000, 8000)))

	car.LeaveSpot()

	parking.ReleaseParkingSpot(spotAvailable)

	car.Leave(spotAvailable)

	// Bloquea el Mutex antes de intentar salir por la puerta
	doorMutex.Lock()

	car.ExitDoor()

	// Desbloquea el Mutex después de salir por la puerta
	doorMutex.Unlock()

	car.GoAway()

	car.Remove()
}

func getRandomNumber(min, max int) float64 {
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	return float64(generator.Intn(max-min+1) + min)
}