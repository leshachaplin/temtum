package master

import (
	"bitbucket.org/dragoninfo/temtum-go-data/repository"
	"bitbucket.org/dragoninfo/temtum-go-master/internal/config"
	"bitbucket.org/dragoninfo/temtum-go-master/internal/dispatcher"
	Kafka "bitbucket.org/dragoninfo/temtum-go-master/internal/kafka"
	"bitbucket.org/dragoninfo/temtum-go-master/internal/operations"
	"bitbucket.org/dragoninfo/temtum-go-master/internal/service"
	websocket "bitbucket.org/dragoninfo/temtum-go-master/internal/websocket"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"
)

var (
	isWatcher bool
)

func Run() {
	isWatcher = true
	log.Info(isWatcher)
	s := make(chan os.Signal)
	signal.Notify(s, os.Interrupt)
	done, cnsl := context.WithCancel(context.Background())
	cfg := config.NewConfig()
	db, err := sqlx.Open(cfg.DbDriver, cfg.ConStr)
	if err != nil {
		log.Fatal("not connect to database", err)
	}

	e := echo.New()

	go func(e *echo.Echo) {
		err := e.Start(fmt.Sprintf(":%d", cfg.ServerPort))
		if err != nil {
			log.Errorf("echo server not start ", err)
		}
	}(e)

	time.Sleep(time.Second * 5)

	kafkaReaderBlock, err := Kafka.New(cfg.TopicTransactions, cfg.TopicBlocks, cfg.KafkaUrl)
	if err != nil {
		log.Errorf("kafka client not connected ", err)
	}
	defer kafkaReaderBlock.Close()

	blocksRepo := repository.NewBlocksRepository(*db)
	transactionRepo := repository.NewTransactionRepository(*db)
	//TODO:status repo
	statusRepo := repository.NewStatusRepository(*db)
	//TODO:uTxO repo
	//uTxoRepo := repository.NewuTxoRepository(*db)

	nodeList, err := statusRepo.GetAllNodes(done)
	if err != nil {
		log.Errorf("node list error ", err)
	}

	srv := websocket.NewServer(done, statusRepo)

	m := service.New(kafkaReaderBlock, blocksRepo, transactionRepo, &isWatcher,
		*cfg, &nodeList, srv)
	d := dispatcher.Dispatcher{Srv: *m}

	srv.Start(e, "/message")

	d.Do(done)
	go d.Srv.ReadBlocks(done, cfg.MaxBytesPerMessage)

	ManagerOfMessages(done, &d)

	<-s
	close(s)
	cnsl()
}

func ManagerOfMessages(ctx context.Context, d *dispatcher.Dispatcher) {
	go func(ctx context.Context) {
		roleTicker := time.NewTicker(time.Second * 40)
		for {
			select {
			case <-roleTicker.C:
				{
					log.Info("tick to change role")
					d.Srv.Operation <- operations.ChangeActiveLieder{}
				}
			case <-ctx.Done():
				{
					return
				}
			}
		}
	}(ctx)
}
