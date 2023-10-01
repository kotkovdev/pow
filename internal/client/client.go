// package client provides client for connection to server.
package client

import (
	"encoding/base64"
	"log/slog"
	"net"
	"strings"

	"github.com/kotkovdev/pow/internal/netutil"
)

const (
	clientProtocol = "tcp"
)

type Client struct {
	conn net.Conn
}

func New(address string) (Client, error) {
	conn, err := net.Dial(clientProtocol, address)
	if err != nil {
		return Client{}, err
	}

	return Client{
		conn: conn,
	}, nil
}

func (c Client) Close() error {
	return c.conn.Close()
}

func (c Client) RequestPuzzle() ([]byte, []byte, error) {
	slog.Info("requesting puzzle")
	if err := c.send(nil); err != nil {
		slog.Error("unable to send request", "error", err)
		return nil, nil, err
	}

	resp, err := c.read()
	if err != nil {
		return nil, nil, err
	}
	slog.Info("got response", "response", resp)

	parts := strings.Split(resp, string(netutil.Separator))
	sourceStr, targetStr := parts[0], parts[1]
	source, err := base64.StdEncoding.DecodeString(sourceStr)
	if err != nil {
		slog.Error("could not parse source hash", "error", err)
		return nil, nil, err
	}
	target, err := base64.StdEncoding.DecodeString(targetStr)
	if err != nil {
		slog.Error("could not parse target hash", "error", err)
		return nil, nil, err
	}

	return source, target, nil
}

func (c Client) ReqeustQuote(payload string) ([]byte, error) {
	if err := c.send([]byte(payload)); err != nil {
		slog.Error("could not send quote request", "error", err, "payload", payload)
		return nil, err
	}

	resp, err := c.read()
	if err != nil {
		slog.Error("could not read response", "error", err)
	}

	return []byte(resp), nil
}

func (c Client) send(body []byte) error {
	return netutil.Send(body, c.conn)
}

func (c Client) read() (string, error) {
	return netutil.Read(c.conn)
}
