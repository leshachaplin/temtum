package websocket

import (
	"bitbucket.org/dragoninfo/temtum-go-master/internal/operations"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"time"
)

func (s *Server) selectOperation(ws *websocket.Conn, op chan operations.Operation, t *time.Ticker, token string) {

	for {
		select {
		case m := <-op:
			{
				switch m.(type) {
				case operations.ExitOperation:
					{
						err := s.exitOperation(m, ws)
						if err != nil {
							log.Errorf(fmt.Sprint(err))
						}
					}
				case operations.SendRoleOperation:
					{
						err := s.sendRoleOperation(m, ws)
						if err != nil {
							log.Errorf(fmt.Sprint(err))
						}
					}
				}
			}
		case <-s.ctx.Done():
			{
				s.wsControllerMutex.Lock()
				delete(s.wsController, token)
				s.wsControllerMutex.Unlock()
				return
			}
		case <-t.C:
			{
				err := s.readHeader(ws)
				if err != nil {
					log.Errorf(fmt.Sprint(err))
				}
			}
		}
	}
}
