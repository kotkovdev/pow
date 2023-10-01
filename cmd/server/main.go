package main

import (
	"flag"
	"log/slog"

	"github.com/kotkovdev/pow/internal/quotes"
	"github.com/kotkovdev/pow/internal/server"
	"github.com/kotkovdev/pow/pkg/challenger"
)

func main() {
	var (
		complexity     int
		dictionaryPath string
		address        string
	)

	flag.IntVar(&complexity, "complexity", 1, "set max complexity of puzzles")
	flag.StringVar(&dictionaryPath, "path", "dictionary.txt", "sets quotes dictionary list path")
	flag.StringVar(&address, "address", ":8080", "sets server address")
	flag.Parse()

	qouteService, err := quotes.New(dictionaryPath)
	if err != nil {
		slog.Error("could not init quotes service", "error", err)
	}

	chal := challenger.NewChallenger(challenger.DefaultSHA256Func, challenger.DefaultSaltGenerateFunc)

	server := server.New(qouteService, chal, complexity)
	if err := server.Serve(address); err != nil {
		slog.Error("could not start server", "error", err)
	}
}
