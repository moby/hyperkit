package u8

import "time"

type Cancellable struct {
	stop chan struct{}
}

func (this *Cancellable) Cancel() {
	this.stop <- struct{}{}
}

func AfterSecond(f func()) Cancellable {
	c := Cancellable{
		stop: make(chan struct{}, 1),
	}
	go func() {
		select {
		case <-time.After(time.Second):
			f()
		case <-c.stop:
		}
	}()

	return c
}
