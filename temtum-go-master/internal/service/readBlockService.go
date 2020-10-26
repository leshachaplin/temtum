package service

import (
	"bitbucket.org/dragoninfo/temtum-go-data/data"
	"bitbucket.org/dragoninfo/temtum-go-master/internal/operations"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/gozstd"
)
//verify block and add to database
func(s *Service) verify(ctx context.Context, message []byte) error {

	block := data.Block{}

	err := json.Unmarshal(message, &block)
	if err != nil {
		log.Errorf(fmt.Sprint("blocks not unmarshal "))
		return err
	}

	decompressData, err := gozstd.Decompress(nil, block.Data)
	_ = decompressData
	transactions, vrf := block.VerifyBlock()

	if vrf {

		err = s.pack(ctx, block, transactions)
		if err != nil {
			return err
		}

	} else {

		err = s.unverifyOperation(block)
		if err != nil {
			return err
		}

	}

	return nil
}

//add block to database
func(s *Service) pack(ctx context.Context, block data.Block, transactions *data.Transactions) error {
	err := s.blocksRepo.Create(ctx, block)
	if err != nil {
		log.Errorf(fmt.Sprint("blocks not added to database"))
		return err
	}

	err = s.transactionRepo.Create(ctx, transactions.T)
	if err != nil {
		log.Errorf(fmt.Sprint("blocks not added to database "))
		return err
	}

	s.lastBlock = block

	return nil
}

//operations which will work if block not verified
func(s *Service) unverifyOperation(block data.Block) error {
	err := s.kafkaReaderBlock.SetOffsetWithConn(block.PreviousOffset, 0)
	if err != nil {
		log.Errorf(fmt.Sprint("offset not set "))
		return err
	}

	err = s.serv.SendToChanel(findActiveLeader(s.nodeList), operations.ExitOperation{Msg: exitOperation})
	if err != nil {
		log.Errorf(fmt.Sprint("exit message not send"))
		return err
	}
	return nil
}
