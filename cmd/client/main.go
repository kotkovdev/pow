package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log/slog"

	"github.com/kotkovdev/pow/pkg/challenger"

	"github.com/kotkovdev/pow/internal/client"
)

func main() {
	var address string
	flag.StringVar(&address, "address", ":8080", "sets server connection address")
	flag.Parse()

	cli, err := client.New(address)
	if err != nil {
		slog.Error("could not connect to server", "error", err, "address", address)
		return
	}
	defer cli.Close()

	source, target, err := cli.RequestPuzzle()
	if err != nil {
		slog.Error("could not get puzzle", "error", err)
		return
	}

	chal := challenger.NewChallenger(challenger.DefaultSHA256Func, challenger.DefaultSaltGenerateFunc)
	answer := chal.SolveRecursive(source, target)

	answerStr := base64.StdEncoding.EncodeToString(answer)
	slog.Info("found solution", "answer", answerStr)

	quote, err := cli.ReqeustQuote(answerStr)
	if err != nil {
		slog.Error("could not request quote", "error", err)
		return
	}

	fmt.Println(string(quote))
}
