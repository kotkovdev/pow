package main

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"log/slog"
	"net"
	"strings"

	"github.com/kotkovdev/pow/internal/util"
	"github.com/kotkovdev/pow/pkg/challenger"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	util.Send(nil, conn)
	slog.Info("sent connection request")

	resp, err := bufio.NewReader(conn).ReadString(byte(util.MessageDelimeter))
	if err != nil {
		slog.Error(err.Error())
		return
	}
	resp = resp[:len(resp)-1]
	slog.Info("got response", "response", resp)
	parts := strings.Split(resp, string(util.Separator))
	sourceStr, targetStr := parts[0], parts[1]
	source, _ := base64.StdEncoding.DecodeString(sourceStr)
	target, _ := base64.StdEncoding.DecodeString(targetStr)

	chal := challenger.NewChallenger(challenger.DefaultSHA256Func)
	answer := chal.SolveRecursive(source, target)
	slog.Info("found solution", "answer", answer, "source", hex.EncodeToString(source), "target", hex.EncodeToString(target))
}
