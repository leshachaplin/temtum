package service

import (
	"bitbucket.org/dragoninfo/temtum-go-data/data"
	"bitbucket.org/dragoninfo/temtum-go-data/repository"
	"bitbucket.org/dragoninfo/temtum-go-leader/internal/config"
	Kafka "bitbucket.org/dragoninfo/temtum-go-leader/internal/kafka"
	"context"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	activeRole    = "active"
	passiveRole   = "passive"
	exitOperation = "exit"
)

type Service struct {
	Operation         chan Operation
	CancelReadBlocks  chan Operation
	CancelWriteBlocks chan Operation
	WsSender          chan Operation
	kafkaReader       *Kafka.Client
	blocksRepo        *repository.BlocksRepository
	transactionRepo   *repository.TransactionRepository
	transactions      data.Transactions
	previousBlock     data.Block
	offset            int64
	config            *config.Config
}

func New(kafkaReader *Kafka.Client, cfg *config.Config, b *repository.BlocksRepository,
	t *repository.TransactionRepository, wsSender chan Operation) *Service {
	return &Service{
		Operation:       make(chan Operation),
		kafkaReader:     kafkaReader,
		WsSender:        wsSender,
		blocksRepo:      b,
		transactionRepo: t,
		previousBlock:   data.Block{},
		config:          cfg,
		transactions:    data.Transactions{T: make([]data.Transaction, 0, 20000)},
	}
}

//Read transactions from kafka
//TODO:Error handler
func (s *Service) ReadTransaction(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case <-s.CancelWriteBlocks:
			{
				log.Info("change active leader")
				return
			}
		default:
			{
				mes, err := s.readTransaction()
				if err != nil {
					log.Error(err)
				}

				log.Info("write block to kafka")

				block, err := s.createBlock(mes)
				if err != nil {
					log.Error(err)
				}

				err = s.packBlock(ctx, block)
				if err != nil {
					log.Error(err)
				}

				err = s.sendBlockToKafka(*block)
				if err != nil {
					log.Error(err)
				}

				err = s.sendLastBlockHeaderToWebsocket(*block)
				if err != nil {
					log.Error(err)
				}

				s.transactions.T = s.transactions.T[:0]

				s.previousBlock = *block

				log.Info("Block successfully send")
			}
		}
	}
}

//Batch read blocks from kafka when leader non-active
func (s *Service) ReadBlocks(ctx context.Context) {
	log.Info("read blocks")
	t := time.NewTicker(time.Second)
	s.readBlock(ctx, t)
}

//operation to change role
func (s *Service) ReadRole(operation Operation) {
	log.Info("read role")

	role := operation.(ReadRoleOperation)
	if role.Role == activeRole {

		s.cancelReadBlocksOperation()

		s.Operation <- ReadTransactionOperation{}
		log.Info("send read operation")

	} else {

		s.cancelWriteBlocksOperation()

		s.Operation <- ReadBlockOperation{}
		log.Info("send write operation")
	}
}
