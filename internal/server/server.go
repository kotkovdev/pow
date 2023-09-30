package server

import (
	"context"
	"log/slog"
	"net"
	"time"

	"github.com/pkg/errors"
)

type server struct{}

const (
	protocol         = "tcp"
	keepAliveTimeout = time.Second
)

func New() server {
	return server{}
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
	slog.Info("handle reques", "local", conn.LocalAddr(), "remote", conn.RemoteAddr())
	slog.Info("handle reques", "connection", conn)
}
