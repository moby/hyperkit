package u8

import (
	"sync"
	"time"
)

type Cancellable struct {
	stop chan struct{}
	wg   *sync.WaitGroup
}

func (this *Cancellable) Cancel() {
	this.stop <- struct{}{}
}

// Wait blocks until the Cancellable goroutine exits.
func (this *Cancellable) Wait() {
	this.wg.Wait()
}

func AfterSecond(f func()) Cancellable {
	c := Cancellable{
		stop: make(chan struct{}, 1),
		wg:   new(sync.WaitGroup),
	}
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		select {
		case <-time.After(time.Second):
			f()
		case <-c.stop:
		}
	}()

	return c
}
