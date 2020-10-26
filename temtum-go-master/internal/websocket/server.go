package websocket

import (
	"bitbucket.org/dragoninfo/temtum-go-data/repository"
	"bitbucket.org/dragoninfo/temtum-go-master/internal/operations"
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	wsControllerMutex sync.Mutex
	ctx               context.Context
	stRepo            repository.StatusRepository
	wsController      map[string]chan operations.Operation
}

//Create new server
func NewServer(ctx context.Context, s *repository.StatusRepository) *Server {
	return &Server{
		wsController: make(map[string]chan operations.Operation),
		ctx:          ctx,
		stRepo:       *s,
	}
}

// New websocket GET request which create new chanel in map where key is token from request
func (s *Server) Start(e *echo.Echo, path string) {

	e.GET(path, func(c echo.Context) error {

		token := c.Request().Header.Get("secret")
		//TODO: verification
		//if s.VerifyToken(token) {
		//	s.wsControllerMutex.Lock()
		//	s.wsController[token] = make(chan operations.Operation)
		//	s.wsControllerMutex.Unlock()
		//}

		//test
		s.wsControllerMutex.Lock()
		s.wsController[token] = make(chan operations.Operation)
		s.wsControllerMutex.Unlock()

		return s.listen(token, c.Response(), c.Request())
	})
}

//verify token
//find token in database
func (s *Server) VerifyToken(tokenRequest string) bool {
	if err := s.stRepo.IfExistToken(s.ctx, tokenRequest); err != nil {
		return false
	}
	return true
}

//sends the operation to the channel that corresponds to the node defined by the token
func (s *Server) SendToChanel(token string, operation operations.Operation) error {
	var err error
	s.wsControllerMutex.Lock()
	if s.wsController[token] == nil {
		err = errors.New("nil pointer exception")
		s.wsControllerMutex.Unlock()
		return err
	}
	s.wsController[token] <- operation
	s.wsControllerMutex.Unlock()
	return nil
}

// websocket handler function which listen headers from all nodes
//we have one connection, but when a request comes to connect to web sockets from a node,
//we verify its token and start listening to requests in a separate goroutine and, if necessary,
//send a signal to the node to complete the operation, we put the signal in the channel
//corresponding to the record in the map for this node
func (s *Server) listen(token string, wr http.ResponseWriter, req *http.Request) error {
	s.wsControllerMutex.Lock()
	op := s.wsController[token]
	s.wsControllerMutex.Unlock()
	t := time.NewTicker(time.Second)
	websocket.Handler(func(ws *websocket.Conn) {
		fmt.Println("handle server connection")
		defer ws.Close()
		s.selectOperation(ws, op, t, token)
	}).ServeHTTP(wr, req)
	return nil
}
