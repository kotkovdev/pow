package util

import (
	"net"
)

const (
	MessageDelimeter = '\n'
	Separator        = '|'
)

func Send(body []byte, conn net.Conn) error {
	_, err := conn.Write(append(body, MessageDelimeter))
	return err
}
