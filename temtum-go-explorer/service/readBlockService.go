package service

import (
	"bitbucket.org/dragoninfo/temtum-go-data/data"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

// read block from kafka,verify and add to database
func(s *Service) readBlock(ctx context.Context, t *time.Ticker) {
	for {
		select {
		case <-s.CancelReadBlocks:
			{
				ctx.Done()
				return
			}
		case <-t.C:
			{
				err := s.kafkaReaderBlock.BatchRead(s.config.MaxBlocksAmount,
					s.config.Timeout,
					s.config.BatchSize,
					s.config.MaxBytesPerMessage,
					func(m []byte) error {

						blocks := &data.Blocks{}

						err := json.Unmarshal(m, blocks)
						if err != nil {
							log.Errorf(fmt.Sprintf("blocks not unmarshal %s", err))
							return err
						}

						err = s.packBlock(ctx, blocks)
						if err != nil {
							log.Errorf(fmt.Sprint(err))
						}

						return nil
					})
				if err != nil {
					log.Errorf(fmt.Sprintf("message not read %s", err))
				}

				log.Info("block successfully synchronized")
			}
		}
	}
}

//verify block and add to DB
func(s *Service) packBlock(ctx context.Context, blocks *data.Blocks) error {
	for i := 0; i < len(blocks.B); i++ {
		transactions, vrf := blocks.B[i].VerifyBlock()
		if vrf {

			err := s.blocksRepo.Create(ctx, blocks.B[i])
			if err != nil {
				log.Errorf(fmt.Sprint("blocks not added to database "))
				return err
			}

			err = s.transactionRepo.Create(ctx, transactions.T)
			if err != nil {
				log.Errorf(fmt.Sprint("transactions not added to database "))
				return err
			}

			err = s.send(blocks.B[i])
			if err != nil {
				return err
			}

		} else {
			break
		}
	}

	return nil
}

//Send header of last readed block to master
func(s *Service) send(block data.Block) error {
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
		log.Errorf(fmt.Sprint("message not Send "))
		return err
	}

	_, err = s.websocketClient.Write(msg)
	if err != nil {
		log.Errorf(fmt.Sprint("message not Send "))
		return err
	}

	return nil
}
