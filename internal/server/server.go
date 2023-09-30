package server

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/hex"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/kotkovdev/pow/pkg/challenger"

	"github.com/kotkovdev/pow/internal/util"
)

type connection struct {
	allowed bool
	expires time.Time
}

type requests *sync.Map

type server struct {
	challenger challenger.Challenger
	requests   requests
}

const (
	protocol         = "tcp"
	keepAliveTimeout = time.Second
)

func New() server {
	return server{
		challenger: challenger.NewChallenger(challenger.DefaultSHA256Func),
		requests:   new(sync.Map),
	}
}

func (s *server) Serve(address string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	slog.Info("starting listener", "address", address)
	listener, err := net.Listen(protocol, address)
	if err != nil {
		return errors.Wrap(ErrAccept, err.Error())
	}
	defer listener.Close()

	// garbage collector.
	go func() {

	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error(errors.Wrap(ErrAccept, err.Error()).Error())
		}

		go func() {
			defer conn.Close()
			s.Handle(ctx, conn)
		}()
	}
}

func (s *server) Handle(ctx context.Context, conn net.Conn) {
	slog.Info("handle request", "local", conn.LocalAddr(), "remote", conn.RemoteAddr())
	body, err := bufio.NewReader(conn).ReadBytes(byte(util.MessageDelimeter))
	if err != nil {
		slog.Error("error read message", "error", err)
	}

	slog.Info("body", "body", body)
	if string(body) == string(util.MessageDelimeter) {
		s.HandleConnection(conn)
		return
	}

	s.HandleSolution(body, conn)
}

func (s *server) HandleConnection(conn net.Conn) {
	puzzle, err := s.challenger.CreatePuzzle([]byte(conn.RemoteAddr().String()), time.Now(), 2)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Info("generated puzzle", "source", hex.EncodeToString(puzzle.Source), "target", hex.EncodeToString(puzzle.Target))
	sourceMsg := []byte(base64.StdEncoding.EncodeToString(puzzle.Source))
	targetMsg := []byte(base64.StdEncoding.EncodeToString(puzzle.Target))
	message := append(append(sourceMsg, util.Separator), targetMsg...)
	if err := util.Send(message, conn); err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Info("sent response", "message", message)
}

func (s *server) HandleSolution(body []byte, conn net.Conn) {
	util.Send([]byte("asdasd"), conn)
}
