package dispatcher

import (
	"bitbucket.org/dragoninfo/temtum-go-master/internal/operations"
	"bitbucket.org/dragoninfo/temtum-go-master/internal/service"
	"context"
)

type Dispatcher struct {
	Srv service.Service
}

func (d *Dispatcher) Do(done context.Context) {
	go func(d *Dispatcher) {
		for {
			select {
			case m := <-d.Srv.Operation:
				switch m.(type) {
				case operations.ChangeActiveLieder:
					{
						go d.Srv.ChangeActiveLieder()
					}
				}
			case <-done.Done():
				{
					return
				}
			}
		}
	}(d)
}
