package explorer

import (
	"bitbucket.org/dragoninfo/temtum-go-data/repository"
	"bitbucket.org/dragoninfo/temtum-go-explorer/config"
	Kafka "bitbucket.org/dragoninfo/temtum-go-explorer/kafka"
	"bitbucket.org/dragoninfo/temtum-go-explorer/service"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

func Run() {
	s := make(chan os.Signal)
	signal.Notify(s, os.Interrupt)
	done, cnsl := context.WithCancel(context.Background())
	cfg := config.NewConfig()

	db, err := sqlx.Open(cfg.DbDriver, cfg.ConStr)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Canno't connect to database %s", err))
	}

	e := echo.New()

	go func(e *echo.Echo) {
		err := e.Start(fmt.Sprintf(":%d", cfg.ServerPort))
		if err != nil {
			log.Errorf(fmt.Sprintf("echo server not start %s", err))
		}
	}(e)

	time.Sleep(time.Second * 5)
	ws, err := connectToWebsocket(cfg.Host, cfg.MsgPath)
	if err != nil {
		log.Errorf(fmt.Sprintf("websocket client not connected %s", err))
	}
	defer ws.Close()

	kafkaReaderBlock, err := Kafka.New(cfg.TopicBlocks, cfg.KafkaUrl, cfg.BatchReadGroup)
	if err != nil {
		log.Errorf(fmt.Sprintf("kafka client not connected %s", err))
	}
	defer kafkaReaderBlock.Close()

	blocksRepo := repository.NewBlocksRepository(*db)
	transactionRepo := repository.NewTransactionRepository(*db)
	//uTxoRepo := repository.NewuTxoRepository(*db)

	l := service.New(kafkaReaderBlock, blocksRepo, transactionRepo, ws, cfg)

	go l.ReadBlocks(done)

	<-s
	close(s)
	cnsl()
}

func connectToWebsocket(host, path string) (*websocket.Conn, error) {
	location := &url.URL{
		Scheme: "ws",
		Host:   host,
		Path:   path,
	}

	origin := &url.URL{
		Scheme: "http",
		Host:   host,
		Path:   "/",
	}

	ws, err := websocket.DialConfig(&websocket.Config{
		Location: location,
		Origin:   origin,
		Version:  websocket.ProtocolVersionHybi13,
		Protocol: []string{},
		Header:   http.Header{"secret": []string{"ExplorerHash"}},
	})
	if err != nil {
		log.Errorf(fmt.Sprint("websocket client not connected"))
		return nil, err
	}
	return ws, nil
}