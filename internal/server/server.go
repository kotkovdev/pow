package server

import (
	"context"
	"encoding/base64"
	"io"
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

type server struct {
	challenger challenger.Challenger
	requests   sync.Map
}

const (
	protocol         = "tcp"
	keepAliveTimeout = time.Second * 5
	requestDeadline  = time.Second * 3
)

func New() server {
	return server{
		challenger: challenger.NewChallenger(challenger.DefaultSHA256Func),
		requests:   sync.Map{},
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
			defer func() {
				conn.Close()
				slog.Info("connection closed")
			}()
			s.Handle(ctx, conn)
		}()
	}
}

func (s *server) Handle(ctx context.Context, conn net.Conn) {
	ctx, cancel := context.WithTimeout(ctx, keepAliveTimeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := s.handle(conn); err != nil {
				if !errors.Is(err, io.EOF) {
					slog.Error(err.Error())
				}
				return
			}
		}
	}
}

func (s *server) handle(conn net.Conn) error {
	// body, err := bufio.NewReader(conn).ReadBytes(byte(util.MessageDelimeter))
	body, err := util.Read(conn)
	if err != nil {
		return err
	}

	slog.Info("handle request", "remote", conn.RemoteAddr(), "body", body)

	if body == "" {
		s.HandleConnection(conn)
	} else {
		s.HandleSolution([]byte(body), conn)
	}

	return nil
}

func (s *server) HandleConnection(conn net.Conn) {
	puzzle, err := s.challenger.CreatePuzzle([]byte(conn.RemoteAddr().String()), time.Now(), 2)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	encoder := base64.StdEncoding

	slog.Info(
		"generated puzzle",
		"source", encoder.EncodeToString(puzzle.Source),
		"target", encoder.EncodeToString(puzzle.Target),
		"original", encoder.EncodeToString(puzzle.Original),
	)

	sourceMsg := []byte(encoder.EncodeToString(puzzle.Source))
	targetMsg := []byte(encoder.EncodeToString(puzzle.Target))
	message := append(append(sourceMsg, util.Separator), targetMsg...)
	if err := util.Send(message, conn); err != nil {
		slog.Error(err.Error())
		return
	}

	payloadHash := encoder.EncodeToString(puzzle.Original)
	s.requests.Store(payloadHash, connection{
		expires: time.Now().Add(requestDeadline),
	})
	slog.Info("store payload hash", "payload", payloadHash)
	slog.Info("sent response", "message", message)
}

func (s *server) HandleSolution(body []byte, conn net.Conn) {
	payloadHash, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		slog.Error("could not decode answer hash", "error", err)
		return
	}
	_ = payloadHash

	req, ok := s.requests.Load(string(body))
	if !ok {
		slog.Error("request not allowed", "payload", string(body))
	}
	_ = req

	if err := util.Send([]byte("response"), conn); err != nil {
		slog.Error("could not send response")
	}
}
