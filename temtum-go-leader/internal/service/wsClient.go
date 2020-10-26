package service

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

func(s *Service) SelectOperation(ctx context.Context, w *websocket.Conn, msg chan Operation) {
	for {
		select {
		case message := <-msg:
			switch message.(type) {
			case SendLastBlockHeaderOperation:
				{
					s.SendLastBlockHeaderFromPassiveLeader(w, message)
				}
			case SendLastBlockOperation:
				{
					s.SendLastBlockHeaderFromActiveLeader(w, message)
				}
			}
		case <-ctx.Done():
			{
				return
			}
		}
	}
}

func(s *Service) SendLastBlockHeaderFromPassiveLeader( w *websocket.Conn, message Operation) {
	log.Info("send last block header operation")
	mes := message.(SendLastBlockHeaderOperation)
	_, err := w.Write(mes.Msg)
	if err != nil {
		log.Errorf(fmt.Sprintf("message not send to master from passive leader %s", err))
	}
}

func(s *Service) SendLastBlockHeaderFromActiveLeader( w *websocket.Conn, message Operation) {
	log.Info("send last block operation")
	mes := message.(SendLastBlockOperation)
	_, err := w.Write(mes.Msg)
	if err != nil {
		log.Errorf(fmt.Sprintf("message not send to master  from active leader %s", err))
	}
}