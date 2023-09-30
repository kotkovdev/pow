package main

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"log/slog"
	"net"
	"strings"

	"github.com/kotkovdev/pow/pkg/challenger"

	"github.com/kotkovdev/pow/internal/util"
)

func main() {
	const address = ":8080"
	conn, err := net.Dial("tcp", address)
	if err != nil {
		slog.Error("could not connect to server", "error", err, "address", address)
		return
	}
	defer conn.Close()

	util.Send(nil, conn)
	slog.Info("sent connection request")

	resp, err := bufio.NewReader(conn).ReadString(util.MessageDelimeter)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// clear delimeter
	resp = resp[:len(resp)-1]
	slog.Info("got response", "response", resp)

	parts := strings.Split(resp, string(util.Separator))
	sourceStr, targetStr := parts[0], parts[1]
	source, err := base64.StdEncoding.DecodeString(sourceStr)
	if err != nil {
		slog.Error("could not parse source hash", "error", err)
		return
	}
	target, err := base64.StdEncoding.DecodeString(targetStr)
	if err != nil {
		slog.Error("could not parse target hash", "error", err)
		return
	}

	chal := challenger.NewChallenger(challenger.DefaultSHA256Func)
	answer := chal.SolveRecursive(source, target)
	slog.Info("found solution", "answer", hex.EncodeToString(answer), "source", hex.EncodeToString(source), "target", hex.EncodeToString(target))

	// conn, err = net.Dial("tcp", address)
	// if err != nil {
	// 	slog.Error("could not connect to server", "error", err, "address", address)
	// 	return
	// }
	// defer conn.Close()

	answerStr := base64.StdEncoding.EncodeToString(answer)
	if err := util.Send([]byte(answerStr), conn); err != nil {
		slog.Error("could not send request", "error", err)
		return
	}
	slog.Info("sent answer request")

	resp, err = bufio.NewReader(conn).ReadString(util.MessageDelimeter)
	if err != nil {
		slog.Error("could not parse response", "error", err)
		return
	}
	slog.Info("got response", "response", resp)
	return
}
