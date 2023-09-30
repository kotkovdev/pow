package util

import (
	"bufio"
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

func Read(conn net.Conn) (string, error) {
	line, err := bufio.NewReader(conn).ReadString(MessageDelimeter)
	if err != nil {
		return "", err
	}
	return line[:len(line)-1], nil
}
