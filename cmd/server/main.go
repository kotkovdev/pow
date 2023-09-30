package main

import (
	"log/slog"

	"github.com/kotkovdev/pow/internal/server"
)

func main() {
	server := server.New()
	if err := server.Serve(":8080"); err != nil {
		slog.Error("could not start server", "error", err)
	}
}
