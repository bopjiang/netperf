package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	listenAddrPort = flag.String("L", "127.0.0.1:8080", "address and ports to export metrics")
)

func main() {
	log.Fatal(http.ListenAndServe(*listenAddrPort, nil))
}
