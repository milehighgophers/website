package main

import (
	"log"
	"os"
	"time"

	"github.com/milehighgophers/website/data"
	"github.com/milehighgophers/website/server"
)

func main() {
	addr := os.Getenv("TARGET")
	if addr == "" {
		addr = "localhost:8080"
	}
	s := data.NewStore(10 * time.Minute)
	go s.Poll()
	log.Fatal(server.Start(addr, s))
}
