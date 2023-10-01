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

	"github.com/kotkovdev/pow/internal/netutil"
	"github.com/kotkovdev/pow/pkg/challenger"
)

type QuotesService interface {
	GetRandomQuote() (string, error)
}

type Challenger interface {
	CreatePuzzle(req []byte, timestamp time.Time, size int) (*challenger.Puzzle, error)
	SolveRecursive(source, target []byte) []byte
}

// connection describes user connection.
type connection struct {
	expires time.Time
}

// Server is a server instance.
type Server struct {
	challenger    Challenger
	quotesService QuotesService
	requests      sync.Map
	complexity    int
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
func New(quotes QuotesService, challenger Challenger, complexity int) Server {
	return Server{
		challenger:    challenger,
		requests:      sync.Map{},
		quotesService: quotes,
		complexity:    complexity,
	}
}

// Serve runs server listener.
func (s *Server) Serve(address string) error {
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
func (s *Server) Handle(ctx context.Context, conn net.Conn) {
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
func (s *Server) handle(conn net.Conn) error {
	body, err := netutil.Read(conn)
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
func (s *Server) HandleConnection(conn net.Conn) {
	puzzle, err := s.challenger.CreatePuzzle([]byte(conn.RemoteAddr().String()), time.Now(), s.complexity)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	encoder := base64.StdEncoding

	sourceStr := encoder.EncodeToString(puzzle.Source)
	targetStr := encoder.EncodeToString(puzzle.Target)
	payloadHash := encoder.EncodeToString(puzzle.Original)

	slog.Info(
		"generated puzzle",
		"source", sourceStr,
		"target", targetStr,
		"original", payloadHash,
	)

	message := append(append([]byte(sourceStr), netutil.Separator), []byte(targetStr)...)
	if err := netutil.Send(message, conn); err != nil {
		slog.Error(err.Error())
		return
	}

	s.requests.Store(payloadHash, connection{
		expires: time.Now().Add(requestDeadline),
	})
	slog.Info("store payload hash", "payload", payloadHash)
	slog.Info("sent response", "message", message)
}

// HandleSolution checks solution answer and sends random phrase.
func (s *Server) HandleSolution(body []byte, conn net.Conn) {
	value, ok := s.requests.Load(string(body))
	if !ok {
		slog.Error("request not allowed", "payload", string(body))
		return
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

	quote, err := s.quotesService.GetRandomQuote()
	if err != nil {
		slog.Error("coudl not get quote", "error", err)
		return
	}

	if err := netutil.Send([]byte(quote), conn); err != nil {
		slog.Error("could not send response", "error", err)
		return
	}
}
