package service

import (
	"bitbucket.org/dragoninfo/temtum-go-data/data"
	"bitbucket.org/dragoninfo/temtum-go-data/repository"
	"bitbucket.org/dragoninfo/temtum-go-master/internal/config"
	Kafka "bitbucket.org/dragoninfo/temtum-go-master/internal/kafka"
	"bitbucket.org/dragoninfo/temtum-go-master/internal/operations"
	"bitbucket.org/dragoninfo/temtum-go-master/internal/websocket"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	Disable         = 0
	Enable          = 1
	Ready           = 2
	ActiveLeader    = 3
	NotSynchronized = 4
	activeRole      = "active"
	passiveRole     = "passive"
	exitOperation   = "exit"
)

type Service struct {
	Operation             chan operations.Operation
	CancelReadBlocks      chan operations.Operation
	CancelWrite           chan operations.Operation
	CancelWriteInsideRole chan operations.Operation
	kafkaReaderBlock      *Kafka.Client
	nodeList              []data.Status
	blocksRepo            *repository.BlocksRepository
	transactionRepo       *repository.TransactionRepository
	lastBlock             data.Block
	isTransactionReader   bool
	config                config.Config
	serv                  *websocket.Server
}

func New(kafkaBlockReader *Kafka.Client, b *repository.BlocksRepository, t *repository.TransactionRepository,
	isReader *bool, cfg config.Config, nList *[]data.Status, serv *websocket.Server) *Service {
	return &Service{
		Operation:           make(chan operations.Operation),
		kafkaReaderBlock:    kafkaBlockReader,
		nodeList:            *nList,
		blocksRepo:          b,
		transactionRepo:     t,
		lastBlock:           data.Block{},
		isTransactionReader: *isReader,
		config:              cfg,
		serv:                serv,
	}
}

//read block from kafka and save it to DB
//TODO: verification and implementation when block not valid
func (s *Service) ReadBlocks(ctx context.Context, maxBytes int) {

		log.Info("read blocks")
		t := time.NewTicker(time.Second)

		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case <-t.C:
				{
					m, err := s.kafkaReaderBlock.ReadWithConn(maxBytes)
					if err != nil {
						log.Errorf(fmt.Sprintf("message not read %s", err))
					}

					err = s.verify(ctx, m.Value)
					if err != nil {
						log.Errorf(fmt.Sprint(err))
					}

					log.Info("block successfully read")
				}
			}
		}
}
