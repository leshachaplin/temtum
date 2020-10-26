package service

import (
	"bitbucket.org/dragoninfo/temtum-go-data/data"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/gozstd"
	"time"
)

func(s *Service) readTransaction() ([]byte, error) {
	log.Info("read transaction")

	m, offset, err := s.kafkaReader.ReadMessage(s.config.BatchSize,
		s.config.MaxBytesPerMessage,
		s.config.MaxTransactionsAmount,
		s.config.ReadTransactionsTimeout)
	if err != nil {
		log.Errorf(fmt.Sprintf("message not read %s", err))
		return nil, err
	}

	s.offset = offset

	trans := &data.Transactions{}

	err = json.Unmarshal(m, trans)
	if err != nil {
		log.Errorf(fmt.Sprintf("transaction not unmarshal %s", err))
		return nil, err
	}

	s.transactions = *trans

	log.Info("Transaction successfully read")

	return m, nil
}

func(s *Service) createBlock(mes []byte, ) (*data.Block, error) {
	block := data.Block{
		Index:          5,
		PreviousHash:   s.previousBlock.Hash,
		BeaconIndex:    5,
		BeaconValue:    "rtsgt",
		Timestamp:      time.Now().Unix(),
		Hash:           "",
		Data:           mes,
		Offset:         s.offset,
		PreviousOffset: s.previousBlock.Offset,
	}

	err := block.CalculateHash()
	if err != nil {
		log.Errorf(fmt.Sprintf("block not create %s", err))
		return nil, err
	}

	//TODO: Verification
	for i := 0; i < len(s.transactions.T); i++ {
		v := s.transactions.T[i].VerifyTransaction()
		if v {
			s.transactions.T[i].BlockHash = block.Hash
		}
	}

	transWithHash, err := json.Marshal(s.transactions)
	if err != nil {
		log.Errorf(fmt.Sprintf("block not marshal %s", err))
		return nil, err
	}

	block.Data = transWithHash

	return &block, nil
}

func(s *Service) packBlock(ctx context.Context, block *data.Block) error {
	err := s.blocksRepo.Create(ctx, *block)
	if err != nil {
		log.Errorf(fmt.Sprintf("blocks not added to database %s", err))
		return err
	}

	block.Data = gozstd.Compress(nil, block.Data)

	log.Info(block.Hash)

	err = s.transactionRepo.Create(ctx, s.transactions.T)
	if err != nil {
		log.Errorf(fmt.Sprintf("transaction not packaging %s", err))
		return err
	}
	return nil
}

func(s *Service) sendBlockToKafka(block data.Block) error {
	blockMsg, err := json.Marshal(block)
	if err != nil {
		log.Errorf(fmt.Sprintf("block not marshal %s", err))
		return err
	}

	err = s.kafkaReader.WriteMessage(blockMsg)
	if err != nil {
		log.Errorf(fmt.Sprintf("block not send to kafka %s", err))
		return err
	}
	return nil
}

func(s *Service) sendLastBlockHeaderToWebsocket(block data.Block) error {
	message := data.Message{
		Index:        block.Index,
		PreviousHash: block.PreviousHash,
		BeaconIndex:  block.BeaconIndex,
		BeaconValue:  block.BeaconValue,
		Timestamp:    block.Timestamp,
		Hash:         block.Hash,
	}

	msg, err := json.Marshal(message)
	if err != nil {
		log.Errorf(fmt.Sprintf("message not Send %s", err))
		return err
	}

	s.WsSender <- SendLastBlockOperation{Msg: msg}

	return nil
}

