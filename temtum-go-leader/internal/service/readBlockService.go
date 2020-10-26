package service

import (
	"bitbucket.org/dragoninfo/temtum-go-data/data"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

func(s *Service) readBlock(ctx context.Context, t *time.Ticker) {
	for {
		select {
		case <-s.CancelReadBlocks:
			{
				return
			}
		case <-ctx.Done():
			{
				return
			}
		case <-t.C:
			{
				err := s.kafkaReader.BatchRead(s.config.MaxBlocksAmount,
					s.config.MaxBytesPerMessage,
					s.config.BatchSize,
					s.config.BatchTimeout,
					func(m []byte) error {

						blocks := &data.Blocks{}

						err := json.Unmarshal(m, blocks)
						if err != nil {
							log.Errorf(fmt.Sprintf("blocks not unmarshal %s", err))
							return err
						}

						//TODO: Verification
						err = s.saveBlocksToDB(ctx, *blocks)
						if err != nil {
							return err
						}

						return nil
					})
				if err != nil {
					log.Errorf(fmt.Sprintf("message not read %s", err))
				}

				log.Info("block successfully read")
			}
		}
	}
}

func(s *Service) saveBlocksToDB(ctx context.Context, blocks data.Blocks) error {
	for i := 0; i < len(blocks.B); i++ {
		transactions, vrf := blocks.B[i].VerifyBlock()
		if vrf {

			err := s.packaging(ctx, blocks, transactions)
			if err != nil {
				continue
			}

			err = s.sendLastReadedBlockHeaderToWebsocket(blocks.B[i])
			if err != nil {
				return err
			}

		} else {
			break
		}
	}

	return nil
}

func(s *Service) sendLastReadedBlockHeaderToWebsocket(block data.Block) error {
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

	s.WsSender <- SendLastBlockHeaderOperation{Msg: msg}

	return nil
}

func(s *Service) packaging(ctx context.Context, blocks data.Blocks, transactions *data.Transactions) error {
	err := s.blocksRepo.BatchCreate(ctx, blocks.B)
	if err != nil {
		log.Errorf(fmt.Sprintf("blocks not added to database %s", err))
		return err
	}

	err = s.transactionRepo.Create(ctx, transactions.T)
	if err != nil {
		log.Errorf(fmt.Sprintf("message not Send %s", err))
		return err
	}

	return nil
}
