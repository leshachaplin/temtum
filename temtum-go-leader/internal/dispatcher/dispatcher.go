package dispatcher

import (
	"bitbucket.org/dragoninfo/temtum-go-leader/internal/service"
	"context"
	"golang.org/x/net/websocket"
)

type Dispatcher struct {
	Srv service.Service
}

func (d *Dispatcher) Do(done context.Context) {
	for {
		select {
		case m := <-d.Srv.Operation:
			switch m.(type) {
			case service.ReadRoleOperation:
				{
					go d.Srv.ReadRole(m)
				}
			case service.ReadTransactionOperation:
				{
					go d.Srv.ReadTransaction(done)
				}
			case service.ReadBlockOperation:
				{
					go d.Srv.ReadBlocks(done)
				}
			}
		case <-done.Done():
			{
				return
			}
		}
	}
}

func (d *Dispatcher) WsClient(ctx context.Context, w *websocket.Conn, msg chan service.Operation) {
	go d.Srv.ReadRoleFromWS(ctx, w)

	d.Srv.SelectOperation(ctx, w, msg)
}
