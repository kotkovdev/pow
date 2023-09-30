package main

import (
	"github.com/kotkovdev/pow/internal/server"
)

func main() {
	server := server.New()
	if err := server.Serve(":8080"); err != nil {
		panic(err)
	}
}
