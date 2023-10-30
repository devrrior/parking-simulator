package models

type Semaphore struct {
	channel chan bool
}

func NewSemaphore(s int) *Semaphore {
	channel := make(chan bool, s)

	for i := 0; i < s; i++ {
		channel <- true
	}

	return &Semaphore{
		channel: channel,
	}
}

func (s *Semaphore) Wait() {
	<-s.channel
}

func (s *Semaphore) Signal() {
	s.channel <- true
}
