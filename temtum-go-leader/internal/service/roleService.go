package service

import 	log "github.com/sirupsen/logrus"

func(s *Service) cancelReadBlocksOperation() {
	if s.CancelReadBlocks != nil {
		s.CancelReadBlocks <- CancelReadBlocks{}
		log.Info("send cancel read")
	} else {
		s.CancelReadBlocks = make(chan Operation)
		s.CancelReadBlocks <- CancelReadBlocks{}
	}
}

func(s *Service) cancelWriteBlocksOperation() {
	if s.CancelWriteBlocks != nil {
		s.CancelWriteBlocks <- CancelWriteBlocks{}
		log.Info("send cancel write")
	} else {
		s.CancelWriteBlocks = make(chan Operation)
		s.CancelWriteBlocks <- CancelWriteBlocks{}
	}
}
