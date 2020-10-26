package service

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"os"
	"time"
)

func(s *Service) ReadRoleFromWS(ctx context.Context, w *websocket.Conn) {
	buff := make([]byte, 1e6)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case <-ticker.C:
			{
				t := time.Now().Add(time.Second * 10)
				w.SetReadDeadline(t)
				n, err := w.Read(buff)
				if err != nil {
					log.Errorf(fmt.Sprintf("message not read from websocket %s", err))
					continue
				}

				trimBuff := buff[:n]

				if string(trimBuff) == exitOperation {
					s.kafkaReader.Close()
					w.Close()
					os.Exit(3)
				}

				if string(trimBuff) == activeRole || string(trimBuff) == passiveRole {
					s.Operation <- ReadRoleOperation{Role: string(trimBuff)}
				}
			}
		}
	}
}
