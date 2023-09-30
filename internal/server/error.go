package server

import "errors"

var (
	ErrListen = errors.New("error listen socket")
	ErrAccept = errors.New("error accept request")
)
