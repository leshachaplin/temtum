package service

import (
	"bitbucket.org/dragoninfo/temtum-go-data/repository"
	"bitbucket.org/dragoninfo/temtum-go-explorer/config"
	Kafka "bitbucket.org/dragoninfo/temtum-go-explorer/kafka"
	"context"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"time"
)

type Service struct {
	Operation             chan Operation
	CancelReadBlocks      chan Operation
	CancelWrite           chan Operation
	CancelWriteInsideRole chan Operation
	kafkaReaderBlock      *Kafka.Client
	websocketClient       *websocket.Conn
	blocksRepo            *repository.BlocksRepository
	transactionRepo       *repository.TransactionRepository
	config                *config.Config
}

func New(kafkaBlockReader *Kafka.Client,
	b *repository.BlocksRepository, t *repository.TransactionRepository,
	w *websocket.Conn, cfg *config.Config) *Service {
	return &Service{
		Operation:        make(chan Operation),
		kafkaReaderBlock: kafkaBlockReader,
		websocketClient:  w,
		blocksRepo:       b,
		transactionRepo:  t,
		config:           cfg,
	}
}

func (s *Service) ReadBlocks(ctx context.Context) {
		log.Info("read blocks")
		t := time.NewTicker(time.Second)
		s.readBlock(ctx, t)
}
