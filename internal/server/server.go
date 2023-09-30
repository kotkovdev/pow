// package server provides TCP server implementation including proof of work functionally.
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

// connection describes user connection.
type connection struct {
	expires time.Time
}

// server is a server instance.
type server struct {
	challenger challenger.Challenger
	requests   sync.Map
}

const (
	// protocol is a server protocol.
	protocol = "tcp"
	// keepAliveTimeout is a connection keep alive time.
	keepAliveTimeout = time.Second * 5
	// requestDeadline is a time after that the request is accepted expired.
	requestDeadline = time.Second * 3
)

// New returns new server instance.
func New() server {
	return server{
		challenger: challenger.NewChallenger(challenger.DefaultSHA256Func),
		requests:   sync.Map{},
	}
}

// Serve runs server listener.
func (s *server) Serve(address string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	slog.Info("starting listener", "address", address)
	listener, err := net.Listen(protocol, address)
	if err != nil {
		return errors.Wrap(ErrAccept, err.Error())
	}
	defer listener.Close()

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

// Handle handles incomming request.
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

// handle handles incomming request.
func (s *server) handle(conn net.Conn) error {
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

// HandleConnection handles incomming request and generates puzzle for solve it on client side.
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

// HandleSolution checks solution answer and sends random phrase.
func (s *server) HandleSolution(body []byte, conn net.Conn) {
	value, ok := s.requests.Load(string(body))
	if !ok {
		slog.Error("request not allowed", "payload", string(body))
	}
	slog.Info("request allowed")
	switch connect := value.(type) {
	case connection:
		defer s.requests.Delete(string(body))
		if connect.expires.Before(time.Now()) {
			slog.Error("connection expired")
			return
		}
	default:
		slog.Error("could not cast connection")
		return
	}

	if err := util.Send([]byte("response"), conn); err != nil {
		slog.Error("could not send response")
	}
}
