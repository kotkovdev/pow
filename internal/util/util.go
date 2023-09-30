// package util provides helping functions.
package util

import (
	"bufio"
	"net"
)

const (
	// MessageDelimeter is a message delimeter.
	MessageDelimeter = '\n'
	// Separator separates message parts.
	Separator = '|'
)

// Send sends message to provided connection.
func Send(body []byte, conn net.Conn) error {
	_, err := conn.Write(append(body, MessageDelimeter))
	return err
}

// Read reads message from connection and removes message delimeter.
func Read(conn net.Conn) (string, error) {
	line, err := bufio.NewReader(conn).ReadString(MessageDelimeter)
	if err != nil {
		return "", err
	}
	return line[:len(line)-1], nil
}
