package websocket

import (
	"bitbucket.org/dragoninfo/temtum-go-master/internal/operations"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"time"
)

func(s *Server) exitOperation(m operations.Operation, ws *websocket.Conn) error {
	log.Info("send exit")
	mes := m.(operations.ExitOperation)
	err := websocket.Message.Send(ws, mes.Msg)
	if err != nil {
		log.Errorf(fmt.Sprint("message not send "))
		return err
	}
	log.Info("end")

	return nil
}

func(s *Server) sendRoleOperation(m operations.Operation, ws *websocket.Conn) error {
	log.Info("send role")
	mes := m.(operations.SendRoleOperation)
	err := websocket.Message.Send(ws, mes.Msg)
	if err != nil {
		log.Errorf(fmt.Sprint("message not send "))
		return err
	}
	log.Info("send role end")

	return nil
}

func(s *Server) readHeader(ws *websocket.Conn) error {
	msg := ""
	log.Info("receive begin")

	t := time.Now().Add(time.Second)
	ws.SetReadDeadline(t)
	err := websocket.Message.Receive(ws, &msg)
	if err != nil {
		log.Errorf(fmt.Sprint("message not read "))
		return err
	}

	log.Info("receive end")

	return nil
}
