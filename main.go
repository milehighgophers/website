package main

import (
	"log"

	"github.com/milehighgophers/website/server"
)

func main() {
	addr := "localhost:8080"
	log.Fatal(server.Start(addr))
}
