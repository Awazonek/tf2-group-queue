package main

import (
	"github.com/awazonek/tf2-group-queue/internal/server"
)

func main() {

	s := server.NewServer()
	s.Start()
}
