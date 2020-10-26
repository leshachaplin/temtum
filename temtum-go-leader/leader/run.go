package leader

import (
	"bitbucket.org/dragoninfo/temtum-go-data/repository"
	"bitbucket.org/dragoninfo/temtum-go-leader/internal/config"
	"bitbucket.org/dragoninfo/temtum-go-leader/internal/dispatcher"
	Kafka "bitbucket.org/dragoninfo/temtum-go-leader/internal/kafka"
	"bitbucket.org/dragoninfo/temtum-go-leader/internal/service"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
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
	msg := make(chan service.Operation)
	signal.Notify(s, os.Interrupt)
	done, cnsl := context.WithCancel(context.Background())
	cfg := config.NewConfig()
	db, err := sqlx.Open(cfg.DbDriver, cfg.ConStr)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Canno't connect to database %s", err))
	}

	e := echo.New()

	go func(e *echo.Echo) {
		err := e.Start(fmt.Sprintf(":%d", cfg.ClientPort))
		if err != nil {
			log.Fatalf(fmt.Sprintf("echo server not start %s", err))
		}
	}(e)

	time.Sleep(time.Second * 5)
	ws, err := connectToWebsocket(cfg.Host, cfg.MsgPath)
	if err != nil {
		log.Errorf(fmt.Sprint(err))
	}
	defer ws.Close()

	time.Sleep(time.Second * 5)

	kafkaReader, err := Kafka.New(cfg.TopicTransactions, cfg.TopicBlocks,
		cfg.KafkaUrl, cfg.BatchBlockGroup, cfg.BlockGroup, cfg.Network)
	if err != nil {
		log.Errorf(fmt.Sprintf("kafka client not connected %s", err))
	}
	defer kafkaReader.Close()

	blocksRepo := repository.NewBlocksRepository(*db)
	transactionRepo := repository.NewTransactionRepository(*db)
	//TODO: UTxO repo
	//uTxoRepo := repository.NewuTxoRepository(*db)

	l := service.New(kafkaReader, cfg, blocksRepo, transactionRepo, msg)

	d := dispatcher.Dispatcher{Srv: *l}

	go d.WsClient(done, ws, msg)

	go d.Do(done)

	d.Srv.Operation <- service.ReadTransactionOperation{}

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
		Header:   http.Header{"secret": []string{"leaderHash"}},
	})
	if err != nil {
		log.Errorf(fmt.Sprint("websocket client not connected"))
		return nil, err
	}
	return ws, nil
}
